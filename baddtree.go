package baddtree

import(
    "fmt"
    //"reflect"
    "unsafe"
    "os"
      )

type node struct {
    //void ** pointers;
    keys []int64
    pointers []unsafe.Pointer
    parent *node
    isLeaf bool
    numKeys int8
    next *node; // Used for queue.
}

type Tree struct {
    root *node
    len int8
}


func Init(len int8) *Tree{
    newTree := new(Tree)
    newTree.root = nil
    newTree.len = len
    return newTree
}

func(newTree *Tree) makeNode() *node{
    newnode := new(node)
    newnode.keys = make([]int64, newTree.len)
    newnode.pointers = make([]unsafe.Pointer, newTree.len)
    newnode.isLeaf = false
    newnode.parent = nil
    newnode.next = nil
    newnode.numKeys = 0
    return newnode
}

func(newTree *Tree) makeLeaf() *node{
    tmp := newTree.makeNode()
    tmp.isLeaf = true
    return tmp
}

func(newTree *Tree) createNewTree(key int64, value unsafe.Pointer){
    newnode := newTree.makeLeaf()
    newnode.keys[0] = key
    newnode.pointers[0] = value
    newnode.pointers[newTree.len - 1] = nil
    newnode.numKeys = 1
    newTree.root = newnode
}



func cut(len int8) int8{
    if len % 2 == 0{
        return len/2
    }
    return len/2 + 1
}

func(newTree *Tree) insertIntoNewRoot(key int64, leftLeaf *node, rightLeaf *node){
    parent := newTree.makeNode()
    parent.keys[0] = key
    parent.pointers[0] = unsafe.Pointer(leftLeaf)
    parent.pointers[1] = unsafe.Pointer(rightLeaf)
    parent.parent = nil
    parent.numKeys++
    leftLeaf.parent = parent
    rightLeaf.parent = parent
    newTree.root = parent
}

func getLeftIndex(parent *node, left *node) int8{
    var i int8
    leftPointer := unsafe.Pointer(left)
    for (i <= parent.numKeys && parent.pointers[i] != leftPointer){
        i++
    }     
    return i
}

func(newTree *Tree) insertIntoParentAfterSplitting(key int64, oldNode *node, right *node, leftIndex int8){
    var i,j int8
    tmpPointers := make([]unsafe.Pointer, newTree.len + 1)
    tmpKeys := make([]int64, newTree.len)
    for i,j=0,0; j < oldNode.numKeys + 1; i,j=i+1,j+1 {
        if i == leftIndex+1 {
            i++
        }
        tmpPointers[i] = oldNode.pointers[j]
    }
    for i,j = 0,0; j < oldNode.numKeys; i,j= i + 1,j + 1 {
        if i == leftIndex{
            i++
        }
        tmpKeys[i] = oldNode.keys[j]
    }
    tmpPointers[leftIndex + 1] = unsafe.Pointer(right)
    tmpKeys[leftIndex] = key

    /*
    for i = 0; i <= oldNode.numKeys; i++ {
        fmt.Print(tmpKeys[i]," ")
    }
    fmt.Println()
    */
    

    split := cut(newTree.len)
    oldNode.numKeys = 0
    for i=0; i < split - 1; i++ {
        oldNode.pointers[i] = tmpPointers[i]
        oldNode.keys[i] = tmpKeys[i]
        oldNode.numKeys++
    }
    oldNode.pointers[i] = tmpPointers[i]
    topKey := tmpKeys[split - 1]

    newNode := newTree.makeNode()
    for i,j = i + 1, 0; i < newTree.len; i,j= i + 1, j + 1 {
        newNode.pointers[j] = tmpPointers[i]
        newNode.keys[j] = tmpKeys[i]
        newNode.numKeys++
    }
    newNode.pointers[j] = tmpPointers[i]
    newNode.parent = oldNode.parent
    for i=0;i<=newNode.numKeys;i=i+1 {
        child := (*node)(unsafe.Pointer(newNode.pointers[i]))
        child.parent = newNode
    }
    newTree.insertIntoParent(topKey, oldNode, newNode)
}

func(newTree *Tree) insertIntoParent(key int64, leftLeaf *node, rightLeaf *node){
    parent := leftLeaf.parent
    if parent == nil {
        newTree.insertIntoNewRoot(key, leftLeaf, rightLeaf)
        return
    }

    leftIndex := getLeftIndex(parent, leftLeaf)
    //fmt.Println("parent.numKeys", parent.numKeys)
    if (parent.numKeys < newTree.len - 1) {
        for i:=parent.numKeys; i > leftIndex; i=i-1 {
            parent.pointers[i + 1] = parent.pointers[i]
            parent.keys[i] = parent.keys[i - 1]
        }
        parent.pointers[leftIndex + 1] = unsafe.Pointer(rightLeaf)
        parent.keys[leftIndex] = key
        parent.numKeys++
        return
    }
    //fmt.Println("parent key : ", key, " leafIndex : ", leftIndex)
    newTree.insertIntoParentAfterSplitting(key, parent, rightLeaf, leftIndex) 
}

func(newTree *Tree) insertIntoLeafAfterSplitting(leaf *node, key int64, value unsafe.Pointer){
    var i, j, insertPosition int8
    insertPosition = 0
    for(insertPosition < newTree.len - 1 && key > leaf.keys[insertPosition]){
        insertPosition++
    }
    tmpNode := newTree.makeLeaf()
    for i,j=0,0; i < leaf.numKeys; i,j=i+1,j+1{
        if j == insertPosition{
            j++
        }
        tmpNode.keys[j] = leaf.keys[i]
        tmpNode.pointers[j] = leaf.pointers[i]
    }
    tmpNode.keys[insertPosition] = key
    tmpNode.pointers[insertPosition] = value

    leaf.numKeys = 0;
    
    split := cut(newTree.len)

    for i = 0; i < split; i=i+1 {
        leaf.keys[i] = tmpNode.keys[i]
        leaf.pointers[i] = tmpNode.pointers[i]
        leaf.numKeys++
    }
    newLeaf := newTree.makeLeaf()
    newLeaf.numKeys = 0
    for i,j = split, 0; i < newTree.len; i,j=i+1,j+1 {
        newLeaf.keys[j] = tmpNode.keys[i]
        newLeaf.pointers[j] = tmpNode.pointers[i]
        newLeaf.numKeys++
    }
    
    for i = leaf.numKeys; i < newTree.len - 1; i++ {
        leaf.pointers[i] = nil
    }

    for i = newLeaf.numKeys; i < newTree.len - 1; i++{
        newLeaf.pointers[i] = nil
    }
    
    newLeaf.pointers[newTree.len - 1] = leaf.pointers[newTree.len - 1]
    newLeaf.parent = leaf.parent
    leaf.pointers[newTree.len - 1] = unsafe.Pointer(newLeaf)
    newKey := newLeaf.keys[0]
    
    //fmt.Println("leaf top key : ", newKey)

    newTree.insertIntoParent(newKey, leaf, newLeaf)
}

func(newTree *Tree) insertIntoLeaf(leaf *node, key int64, value unsafe.Pointer) *node{
    var i, insertPosition int8
    insertPosition = 0
    for(insertPosition < leaf.numKeys && key > leaf.keys[insertPosition]){
        insertPosition++
    }
    for i = leaf.numKeys; i > insertPosition; i-- {
        leaf.keys[i] = leaf.keys[i-1]
        leaf.pointers[i] = leaf.pointers[i-1]
    }
    leaf.keys[insertPosition] = key
    leaf.pointers[insertPosition] = value
    leaf.numKeys++
    return leaf
}

func(newTree *Tree) findLeaf(key int64) *node{
    var i int8
    myNode := newTree.root
    if myNode == nil{
        return nil
    }
    for myNode.isLeaf == false{
        i = 0
        for i < myNode.numKeys{
            if key >= myNode.keys[i]{
                i++
            } else {
                break
            }
        }
        myNode = (*node)(unsafe.Pointer(myNode.pointers[i]))
        
            if myNode == nil {
                //newTree.PrintTree()
            }
    }
    //fmt.Println("find node  numKeys[0]", myNode.keys[0], " find key : ", key)
    return myNode
}

func(newTree *Tree) findPointer(key int64) unsafe.Pointer{
    leaf := newTree.findLeaf(key)
    if leaf == nil {
        return nil
    }
    var i int8
    for i=0; i<leaf.numKeys; i++{
        if leaf.keys[i] == key {
            return leaf.pointers[i]
        }
    }
    return nil
}

func(newTree *Tree) Find(key int64) string{
    leaf := newTree.findLeaf(key)
    if leaf == nil {
        return ""
    }
    var i int8
    for i=0; i<leaf.numKeys; i++{
        if leaf.keys[i] == key {
            value := (*string)(unsafe.Pointer(leaf.pointers[i]))
            return *value
        }
    }
    return ""
}
func(newTree *Tree) Inster(key int64, value string) bool{
    //defer newTree.PrintTree()
    if newTree.root == nil {
        newTree.createNewTree(key, unsafe.Pointer(&value))
        return true
    }
    if newTree.Find(key) != "" {
        return false
    }
    leaf := newTree.findLeaf(key)
    if(leaf.numKeys < newTree.len - 1){
        newTree.insertIntoLeaf(leaf, key, unsafe.Pointer(&value))
        return true
    }
    newTree.insertIntoLeafAfterSplitting(leaf, key, unsafe.Pointer(&value))
    return true
}

func(newTree *Tree) removeEntryFromNode(key int64, delPointer unsafe.Pointer, delNode *node){
    var i, pointersNum int8
    for(key != delNode.keys[i]) {
        i++
    }
    for i=i+1; i < delNode.numKeys; i++{
        delNode.keys[i - 1] = delNode.keys[i]
    }
    i = 0
    if delNode.isLeaf {
        pointersNum = delNode.numKeys
    } else {
        pointersNum = delNode.numKeys + 1
    }
    for(delPointer != delNode.pointers[i]) {
        i++
    }
    for i = i + 1; i < pointersNum; i++ {
        delNode.pointers[i - 1] = delNode.pointers[i]
    }
    delNode.numKeys--
}

func(newTree *Tree) adjustRoot(){
    root := newTree.root
    if root.numKeys > 0 {
        return
    }
    if(root.isLeaf == false) {
        newRoot := (*node)(unsafe.Pointer(root.pointers[0]))
        newTree.root = newRoot
        newRoot.parent = nil
        return
    }
    root = nil
    return
}

func(newTree *Tree) getNeighborIndex(delNode *node) int8{
    var i int8
    parent := delNode.parent
    pointer := unsafe.Pointer(delNode)
    for i=0;i<parent.numKeys;i++ {
        if parent.pointers[i] == pointer {
            break
        }
    }
    return i-1
}

func(newTree *Tree) coalesceNodes(primeKey int64, delNode *node, neighborNode *node, neighborIndex int8){
    var i,j int8
    
    if neighborIndex == -1 {
        tmp := delNode
        delNode = neighborNode
        neighborNode = tmp
    }
    
    neighborInsertIndex := neighborNode.numKeys

    if delNode.isLeaf == false {
        neighborNode.keys[neighborInsertIndex] = primeKey
        neighborNode.numKeys++
        end := delNode.numKeys
        for i,j=neighborInsertIndex + 1, 0; j < end; i,j=i+1,j+1 {
            neighborNode.keys[i] = delNode.keys[j]
            neighborNode.pointers[i] = delNode.pointers[j]
            neighborNode.numKeys++
            delNode.numKeys--
        }
        neighborNode.pointers[i] = delNode.pointers[j]
        for i=0; i <= neighborNode.numKeys; i++ {
            tmp := (*node)(unsafe.Pointer(neighborNode.pointers[i]))
            tmp.parent = neighborNode
        }
    }else{
        for i,j=neighborNode.numKeys,0; j<delNode.numKeys; i,j=i+1,j+1 {
            neighborNode.keys[i] = delNode.keys[j]
            neighborNode.pointers[i] = delNode.pointers[j]
            neighborNode.numKeys++
        }
        neighborNode.pointers[newTree.len - 1] = delNode.pointers[newTree.len - 1]
    }
    /*
    fmt.Println("\n ",unsafe.Pointer(&delNode))
        for i=0;i<=delNode.parent.numKeys;i++ {
            fmt.Println(delNode.parent.pointers[i])
        }
    fmt.Println("\n ",unsafe.Pointer(delNode))
    */
    newTree.deleteEntry(primeKey, unsafe.Pointer(delNode), delNode.parent)
}

func(newTree *Tree) redistributeNodes(delNode *node, neighborNode *node, neighborIndex int8, primeKeyIndex int8, primeKey int64){
    var i int8
    if neighborIndex == -1{
        if (delNode.isLeaf == true) {
            delNode.keys[delNode.numKeys] = neighborNode.keys[0]
            delNode.pointers[delNode.numKeys] = neighborNode.pointers[0]
            neighborNode.parent.keys[primeKeyIndex] = neighborNode.keys[1]           
        } else {
            delNode.keys[delNode.numKeys] = primeKey
            delNode.pointers[delNode.numKeys + 1] = neighborNode.pointers[0]
            tmp := (*node)(unsafe.Pointer(neighborNode.pointers[0]))
            tmp.parent = delNode
            delNode.parent.keys[primeKeyIndex] = neighborNode.keys[0]
        }
        for i=0; i<neighborNode.numKeys; i++ {
            neighborNode.keys[i] = neighborNode.keys[i + 1]
            neighborNode.pointers[i] = neighborNode.pointers[i + 1]
        }
        if delNode.isLeaf == false {
            neighborNode.pointers[i] = neighborNode.pointers[i+1]
        }
    } else {
        if (delNode.isLeaf == false) {
            delNode.pointers[delNode.numKeys + 1] = delNode.pointers[delNode.numKeys]
        }
        for i=delNode.numKeys; i>0;i-- {
            delNode.pointers[i] = delNode.pointers[i - 1]
            delNode.keys[i] = delNode.keys[i - 1]
        }
        if (delNode.isLeaf) {
            delNode.pointers[0] = neighborNode.pointers[neighborNode.numKeys - 1]
            delNode.keys[0] = neighborNode.keys[neighborNode.numKeys - 1]
            neighborNode.pointers[neighborNode.numKeys - 1] = nil
            delNode.parent.keys[primeKeyIndex] = delNode.keys[0]
        } else {
            delNode.pointers[0] = neighborNode.pointers[neighborNode.numKeys]
            //delNode.keys[0] = neighborNode.parent.keys[primeKeyIndex]
            delNode.keys[0] = primeKey
            neighborNode.pointers[neighborNode.numKeys] = nil
            delNode.parent.keys[primeKeyIndex] = neighborNode.keys[neighborNode.numKeys - 1]
            tmp := (*node)(unsafe.Pointer(delNode.pointers[0]))
            tmp.parent = delNode
        }
    }
    delNode.numKeys++
    neighborNode.numKeys--
}

func(newTree *Tree) deleteEntry(key int64, value unsafe.Pointer, delNode *node){
    var minKeys,primeKeyIndex,neighborNodeIndex,capacity int8
    newTree.removeEntryFromNode(key,value, delNode)
    if delNode == newTree.root {
        newTree.adjustRoot()
        return
    }
    if delNode.isLeaf {
        minKeys = cut(newTree.len - 1)
    } else {
        minKeys = cut(newTree.len) - 1
    }
    if delNode.numKeys >= minKeys {
        return
    }
    neighborIndex := newTree.getNeighborIndex(delNode)
    if neighborIndex == -1 {
        primeKeyIndex = 0
    } else {
        primeKeyIndex = neighborIndex
    }
    primeKey := delNode.parent.keys[primeKeyIndex]
    if neighborIndex == -1 {
        neighborNodeIndex = 1
    } else {
        neighborNodeIndex = neighborIndex
    }
    neighborNode := (*node)(unsafe.Pointer(delNode.parent.pointers[neighborNodeIndex]))

    if delNode.isLeaf {
        capacity = newTree.len
    } else {
        capacity = newTree.len - 1
    }
    if (neighborNode.numKeys + delNode.numKeys < capacity){
        newTree.coalesceNodes(primeKey, delNode, neighborNode, neighborIndex)
    } else {
        newTree.redistributeNodes(delNode, neighborNode, neighborIndex, primeKeyIndex, primeKey)
    }
}

func(newTree *Tree) Delete(key int64){
    valuePointer := newTree.findPointer(key)
    leaf := newTree.findLeaf(key)
    if (leaf != nil && valuePointer != nil) {
       newTree.deleteEntry(key, valuePointer, leaf) 
    }
}

func Test(a int) {
    fmt.Println("haha")
}

func(newTree *Tree) PrintTree(){
    var list []*node
    var i, now, next int8

    var flag bool
    flag = false
    list = append(list, newTree.root)
    //parent := newTree.root
    now++
    for len(list) > 0 {
        /*
        if list[0].parent != nil && parent != list[0].parent {
            fmt.Println()
        }
        */
        for i=0;i<list[0].numKeys;i++ {
            fmt.Print(list[0].keys[i], " ")
                if (list[0].keys[i] == 0) {
                    flag = true
                }
            /*
            cc := list[0].parent
            if cc == nil {
                break
            }
            fmt.Print(list[0].keys[i],"(",cc.keys[0],") ")
            */
        }
        fmt.Print("|")
        if list[0].isLeaf == false {
            for i=0;i<=list[0].numKeys;i++ {
                next++
                list = append(list, (*node)(unsafe.Pointer(list[0].pointers[i])))
            }
        }
        //parent = list[0].parent
        list = list[1:]
        now--
        if now == 0 {
            fmt.Println()
            now = next
            next = 0
        }
    }
    fmt.Println("\n--------------------------------")
        if flag {
            os.Exit(1)
        }
}
