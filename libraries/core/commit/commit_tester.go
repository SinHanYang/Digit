package main

import "fmt"

func main(){
    var hash Hash
    hash.hash="1"
    new_commit("2020","harry",hash,"abc")
    var ans=find_commit(hash)
    fmt.Printf("ans:",ans.author)
}