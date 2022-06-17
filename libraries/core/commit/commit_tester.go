package main

import(
     "fmt"
     "time"
)

func main(){
    var hash Hash
    hash.hash_value="hash"
    t:=time.Now()
    new_commit(t,"harry",hash,"abc")
    var ans=find_commit(hash)
    fmt.Printf("ans:",ans.time)
}