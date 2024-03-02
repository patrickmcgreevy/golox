package bytecode

import (
	"fmt"
	"strings"
)

var loadFactor float64 = 0.5

// Hash function interface
type Hasher interface {
	Hash(LoxString) int
}

type HashFunction func(LoxString) int

func (f HashFunction) Hash(s LoxString) int {
	return f(s)
}

func fvnHash(s LoxString) int {
    var hash uint32 = 2166136261
    var prime uint32 = 16777619

    for _, v := range s {
        hash =  hash * prime
        hash = uint32(uint8(v)) ^ hash
    }

    return int(hash)
}

var FVNHashFunction = HashFunction(fvnHash)

type keyValPair struct {
	key LoxString
	val Value
}

type LinearProbingHashMap struct {
	buckets      []keyValPair
	loadFactor   float64
	hashFunction Hasher
}

func NewLinearProbingHashMap() LinearProbingHashMap {
	return LinearProbingHashMap{
		buckets:      make([]keyValPair, 1000),
		loadFactor:   0,
		hashFunction: FVNHashFunction,
	}
}

func (hashMap LinearProbingHashMap) getIndex(s LoxString) int {
    return FVNHashFunction.Hash(s) % cap(hashMap.buckets)
}

func (hashMap LinearProbingHashMap) String() string {
	str := strings.Builder{}
	str.WriteString("{\n")
	for _, v := range hashMap.buckets {
		if v.key != "" {
			str.WriteString(fmt.Sprintf("\t%s: %v,\n", v.key, v.val))
		}
	}
	str.WriteString("}\n")

	return str.String()
}

func (hashMap *LinearProbingHashMap) Insert(s LoxString, v Value) {
    i := hashMap.getIndex(s)
	for true {
		if hashMap.buckets[i].key == "" || hashMap.buckets[i].key == s {
			hashMap.buckets[i] = keyValPair{s, v}
			break
		} else {
			i++
			if i >= cap(hashMap.buckets) {
				i = 0
			}
		}
	}

	hashMap.loadFactor += 1 / float64(cap(hashMap.buckets))
	// rehash
	if hashMap.loadFactor >= loadFactor {
		oldBuckets := hashMap.buckets
		hashMap.buckets = make([]keyValPair, cap(hashMap.buckets)*2)
		hashMap.loadFactor = 0
		for _, v := range oldBuckets {
			if v.key != "" {
				hashMap.Insert(v.key, v.val)
			}
		}
	}
}

func (hashMap *LinearProbingHashMap) Get(s LoxString) (Value, error) {
    initial_i := hashMap.getIndex(s)
	i := initial_i
	for true {
		pair := hashMap.buckets[i]
		if pair.key == s {
			return pair.val, nil
		}
		i++
		if i >= cap(hashMap.buckets) {
			i = 0
		}
		if i == initial_i {
			return nil, fmt.Errorf("%s is not in the map", s)
		}
	}

	return nil, nil // unreachable
}

func (hashMap *LinearProbingHashMap) Delete(s LoxString) {
    initial_i := hashMap.getIndex(s)
	i := initial_i
	for true {
		pair := hashMap.buckets[i]
		if pair.key == s {
            hashMap.buckets[i].key = ""
            hashMap.loadFactor -= 1 / float64(cap(hashMap.buckets))
            return
		}
		i++
		if i >= cap(hashMap.buckets) {
			i = 0
		}
		if i == initial_i {
            return
		}
	}
}
