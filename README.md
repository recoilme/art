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


### Status
WIP

### Storage format
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
### Benchmarks (Art vs HashMap)

```
$ go test -bench=. -benchmem -count=3 -timeout=1m  > x.txt
$ benchstat x.txt
```
Note: in this benchmark keys are:
 - ints - (in bigendian encodings, many common bytes)
 - words (more realistic)

```
name               time/op
SetArt-8            115ns ±16%
SetHashMap-8        277ns ±10%
GetArt-8           48.5ns ± 1%
AscendArt-8         182ns ± 8%
DescendArt-8        401ns ± 3%
ScanArt-8           176ns ± 4%
GetHashMap-8        116ns ± 8%
GetWordsArt-8       155ns ± 7%
GetWordsHashMap-8   116ns ± 7%

name               alloc/op
SetArt-8            90.0B ± 0%
SetHashMap-8        8.00B ± 0%
GetArt-8            0.00B     
AscendArt-8         0.00B     
DescendArt-8         263B ± 2%
ScanArt-8           0.00B     
GetHashMap-8        0.00B     
GetWordsArt-8       0.00B     
GetWordsHashMap-8   0.00B     

name               allocs/op
SetArt-8             1.00 ± 0%
SetHashMap-8         1.00 ± 0%
GetArt-8             0.00     
AscendArt-8          0.00     
DescendArt-8         0.00     
ScanArt-8            0.00     
GetHashMap-8         0.00     
GetWordsArt-8        0.00     
GetWordsHashMap-8    0.00  
```

[More benchmarks](https://github.com/recoilme/bench_sortedsets)

### Usage

see tests

### Credits

 - [The Adaptive Radix Tree: ARTful Indexing for Main-Memory Databases (Specification)](https://db.in.tum.de/~leis/papers/ART.pdf)

### Author

[Kulibaba Vadim](https://github.com/recoilme)