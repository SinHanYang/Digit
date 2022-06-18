package main

import(
     "fmt"
     "time"
)

func main(){
    var hash ="hash2"
    t2:=time.Now()
    time.Sleep(time.Second)
    t:=time.Now()
    var hash2="hash"
    new_commit(t2,"harry",hash2,"value","second")
    new_commit(t,"head_author",hash,"value","first")
    var ans=find_commit(hash)
    fmt.Printf("ans:",ans.time)
    fmt.Printf("\n")
    var commit=find_head()
    fmt.Printf("head:",commit.author)
    fmt.Printf("\n")
    var list = listcommit()
    for _,s :=range list{
        fmt.Printf(s.message)
        fmt.Printf("\n")
    }
}