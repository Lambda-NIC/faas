package handlers

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/client"
)

//FunctionDeployment contains replicas for each functions
type FunctionDeployment struct {
	SmartNIC    string
	NumReplicas uint64
}

//FunctionDeploymentMeta is a local cache for storing function deployments
type FunctionDeploymentMeta struct {
	LastRefresh        time.Time
	FunctionDeployment []FunctionDeployment
}

// Expired find out whether the cache item has expired with
// the given expiry duration from when it was stored.
func (fm *FunctionDeploymentMeta) Expired(expiry time.Duration) bool {
	return time.Now().After(fm.LastRefresh.Add(expiry))
}

// FunctionCache provides a cache of Function replica counts
type FunctionCache struct {
	Cache  map[string]*FunctionDeploymentMeta
	Expiry time.Duration
	Sync   sync.Mutex
}

// Set replica count for functionName
func (fc *FunctionCache) Set(functionName string,
	functionDeployment []FunctionDeployment) {
	fc.Sync.Lock()
	defer fc.Sync.Unlock()

	if _, exists := fc.Cache[functionName]; !exists {
		fc.Cache[functionName] = &FunctionDeploymentMeta{}
	}

	entry := fc.Cache[functionName]
	entry.LastRefresh = time.Now()
	entry.FunctionDeployment = functionDeployment

}

// Get replica count for functionName
func (fc *FunctionCache) Get(functionName string) (FunctionDeployment, bool) {

	fc.Sync.Lock()
	defer fc.Sync.Unlock()

	replicas := FunctionDeployment{}

	hit := false
	/*
		if val, exists := fc.Cache[functionName]; exists {
			replicas = val.FunctionDeployment
			hit = !val.Expired(fc.Expiry)
		}
	*/
	return replicas, hit
}

/*
var Cache FunctionCache = FunctionCache{
	Cache:  make(map[string]*FunctionDeploymentMeta),
	Expiry: config.CacheExpiry,
}
*/

// CreateEtcdClient creates a client for ETCD deployment
func CreateEtcdClient(etcdMasterIP string, etcdPort string) client.KeysAPI {
	cfg := client.Config{
		Endpoints: []string{fmt.Sprintf("http://%s:%s", etcdMasterIP, etcdPort)},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when
		// the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal("Could not connect to ETCD: " + err.Error())
	}
	kapi := client.NewKeysAPI(c)
	return kapi
}

// CreateDepKey creates a key for deployment
func CreateDepKey(smartNIC string, funcName string) string {
	return fmt.Sprintf("/deployments/smartnic/%s/%s", smartNIC,
		funcName)
}

// CreateFuncKey creates a key for function
func CreateFuncKey(funcName string) string {
	return fmt.Sprintf("/functions/%s", funcName)
}

// EtcdFunctionExists checks if the function exists in etcd.
func EtcdFunctionExists(keysAPI client.KeysAPI, functionName string) bool {
	_, err := keysAPI.Get(context.Background(),
		fmt.Sprintf("/functions/%s", functionName),
		nil)
	// Did not find the function.
	return err == nil
}

// GetSmartNICS returns the list of SmartNICs from ETCD.
func GetSmartNICS(keysAPI client.KeysAPI) []string {
	resp, err := keysAPI.Get(context.Background(), "/smartnics", nil)
	// No smartnics found in deployment.
	if err != nil {
		log.Println("Could not retrieve SmartNICs")
		return nil
	}
	var smartNICs []string
	sort.Sort(resp.Node.Nodes)
	for _, n := range resp.Node.Nodes {
		smartNICs = append(smartNICs, strings.Split(n.Key, "/")[2])
	}
	return smartNICs
}

// GetFunctions returns the list of functions
func GetFunctions(keysAPI client.KeysAPI) ([]string, error) {
	resp, err := keysAPI.Get(context.Background(), "/functions", nil)
	if err != nil {
		log.Println("Could not retrieve functions")
		return nil, err
	}
	var functions []string
	sort.Sort(resp.Node.Nodes)
	for _, n := range resp.Node.Nodes {
		functions = append(functions, strings.Split(n.Key, "/")[2])
	}
	return functions, nil
}

// GetNumDeployments gives the number of deployments for the function.
func GetNumDeployments(keysAPI client.KeysAPI,
	funcName string) uint64 {

	var numReplicas uint64
	smartNICs := GetSmartNICS(keysAPI)
	if smartNICs == nil {
		return 0
	}
	for _, smartNIC := range smartNICs {
		depVal, depErr := keysAPI.Get(context.Background(),
			CreateDepKey(smartNIC, funcName), nil)
		if depErr != nil {
			// Deployment doesn't exist
			continue
		} else {
			numDeps, numDepErr := strconv.ParseUint(depVal.Node.Value, 10, 64)
			if numDepErr == nil {
				numReplicas += numDeps
			}
		}
	}
	return numReplicas
}
