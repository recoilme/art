# Art 
Adaptive Radix Tree done right

## simple & perfomant adaptive radix tree implementation

### Description
Unlike other implementations in Go/C this version store only differents parts of keys. Which leads to dramatically reducing of memory usage in case of storing keys with repeatable fragments

### Advantages

 - Optimize memory usage, correct implementation of compact/prefix tree
 - Minimum allocations on Set/Get (GC friendly)
 - Perfomant
 - Maintains the data in sorted order, which enables additional operations like range scan and prefix lookup

### Disadvantages:
 
 - Binary comparator


## Status: WIP

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
        http://:
         example.com/tag/:
          10:a
          20:b
         some.com:c
```
In binary format:
```
        &{key:[104 116 116 112 58 47 47] val:[] children:[0xc00005f1d0 0xc00005f220 <nil> <nil>] size:2}
         &{key:[101 120 97 109 112 108 101 46 99 111 109 47 116 97 103 47] val:[] children:[0xc00005f130 0xc00005f180 <nil> <nil>] size:2}
          &{key:[49 48] val:[97] children:[] size:0}
          &{key:[50 48] val:[98] children:[] size:0}
         &{key:[115 111 109 101 46 99 111 109] val:[99] children:[] size:0}
```

## Usage

see tests

## Credits

 - [The Adaptive Radix Tree: ARTful Indexing for Main-Memory Databases (Specification)](https://db.in.tum.de/~leis/papers/ART.pdf)