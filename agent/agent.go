package main

import (
	"encoding/json"
	"math/rand"
	"os"

	"github.com/Gimulator-Games/xo-random-agent/world"
	"github.com/Gimulator/client-go"
)

var ch = make(chan client.Object, 100)

func main() {
	cl, err := client.NewClient(ch)
	user := os.Getenv("CLIENT_ID")
	if err != nil {
		panic(err)
	}
	keyRegister := client.Key{
		Type:      "register",
		Namespace: user,
	}
	err = cl.Set(keyRegister, nil)
	if err != nil {
		panic(err)
	}
	keyAction := client.Key{
		Type:      "action",
		Namespace: "",
		Name:      user,
	}
	keyWorld := client.Key{
		Type: "world",
	}
	var w world.World
	cl.Watch(keyWorld)
	for {
		tmp := <-ch
		data := tmp.Value
		json.Unmarshal([]byte(data.(string)), &w)
		if w.Turn == user && w.Result != "end" {
			action := nextMove(w)
			jsonAction, _ := json.Marshal(action)
			cl.Set(keyAction, string(jsonAction))
		}
	}
}

func nextMove(w world.World) world.Action {
	moves := w.Moves
	var availableMoves = make([]int, 0, 9)
	var flag bool
	for i := 0; i < 9; i++ {
		flag = true
		for _, j := range moves {
			if i == j.Pos {
				flag = false
			}
		}
		if flag {
			availableMoves = append(availableMoves, i)
		}
	}
	if len(availableMoves) > 1 {
		rnd := rand.Intn(len(availableMoves) - 1)
		action := world.Action{Pos: availableMoves[rnd]}
		// fmt.Println(action)
		return action
	}
	action := world.Action{Pos: availableMoves[0]}
	// fmt.Println(action)
	return action
}
