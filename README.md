# Art 
Adaptive Radix Tree done right

## Simple & perfomant adaptive radix tree implementation

### Description
Unlike other implementations in Go/C this version store only different parts of keys. Which leads to dramatically reducing of memory usage in case of storing keys with repeatable fragments

### Advantages

 - Optimized memory usage, correct implementation of compact/prefix tree
 - Minimized allocations on Set/Get (GC friendly)
 - Perfomant
 - Store the data in sorted order

### Disadvantages
 
 - Binary comparator only


## Status
WIP

## Storage format
In this example:
```
	tree := art.New()

	tree.Set([]byte("http://example.com/tag/10"), []byte("a"))
	tree.Set([]byte("http://example.com/tag/20"), []byte("b"))
	tree.Set([]byte("http://some.com"), []byte("c"))

	t.Log(tree.StringKeys(true))
```
Tree will looks like:
```
        key:http:// val:
         key:example.com/tag/ val:
          key:10 val:a
          key:20 val:b
         key:some.com val:c
```
In binary format:
```
        &{key:[104 116 116 112 58 47 47] val:[] children:[0xc00005f1d0 0xc00005f220 <nil> <nil>] size:2}
         &{key:[101 120 97 109 112 108 101 46 99 111 109 47 116 97 103 47] val:[] children:[0xc00005f130 0xc00005f180 <nil> <nil>] size:2}
          &{key:[49 48] val:[97] children:[] size:0}
          &{key:[50 48] val:[98] children:[] size:0}
         &{key:[115 111 109 101 46 99 111 109] val:[99] children:[] size:0}
```

## Benchmarks (Art vs HashMap)

```
$ go test -bench=. -benchmem -count=3 -timeout=1m  > x.txt
$ benchstat x.txt
```
Note: in this benchmark keys has many common parts (ints - in bigendian encodings && words with common prefixes)

```
name               time/op
SetArt-8            137ns ±12%
SetHashMap-8        212ns ±10%
GetArt-8           36.8ns ± 4%
GetHashMap-8        115ns ± 2%
GetWordsArt-8       124ns ± 5%
GetWordsHashMap-8  89.6ns ± 4%

name               alloc/op
SetArt-8            90.0B ± 0%
SetHashMap-8        8.00B ± 0%
GetArt-8            0.00B     
GetHashMap-8        0.00B     
GetWordsArt-8       0.00B     
GetWordsHashMap-8   0.00B 
```

## Usage

see tests

## Credits

 - [The Adaptive Radix Tree: ARTful Indexing for Main-Memory Databases (Specification)](https://db.in.tum.de/~leis/papers/ART.pdf)