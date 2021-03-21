package exps

func Hget() {
	key := "test_hget"
	type Petime struct {
		Id   int
		Name string
	}
	petimes := []Petime{
		{1, "M"},
		{2, "S"},
	}
	for _, pt := range petimes {
		rdbHandler.Hset(key, pt.Id, pt.Name)
	}
}
