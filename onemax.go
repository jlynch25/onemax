package main

import (
	"math/rand"
	"sort"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type Entry struct {
	Value string
	Score int
}
type ByScore []Entry

var entries []Entry = []Entry{}

var maxFitness int = 20

func main() {
	population := 100
	totalFitness := 0.0
	var averageFitness []float64
	var averageFitnessIndex []float64

	rand.Seed(time.Now().UnixNano())

	for i := 1; i < population; i++ {
		e := Entry{Value: RandSeq(maxFitness)}
		e.Score = calScore(e)
		entries = append(entries, e)
	}

	//// (iii) just make sure string is in list
	e := Entry{Value: "00000000000000000000"}
	e.Score = calScore(e)
	entries = append(entries, e)

	for j := 1; j <= population; j++ {

		sort.Sort(sort.Reverse(ByScore(entries)))

		totalFitness = 0
		for _, e := range entries {
			totalFitness = totalFitness + float64(e.Score)
		}
		average := totalFitness / float64(len(entries)*maxFitness)
		averageFitness = append(averageFitness, average)
		averageFitnessIndex = append(averageFitnessIndex, float64(j))

		if average == 1 {
			break
		}

		usedPopulationSize := int(float64(len(entries)) * 0.05)

		for i := 1; i < usedPopulationSize; i++ {
			rand1 := rand.Intn(usedPopulationSize)
			rand2 := rand.Intn(usedPopulationSize)
			entries[rand1], entries[rand2] = crossover(entries[rand1], entries[rand2])
			rand3 := rand.Intn(usedPopulationSize)
			entries[rand3] = mutation(entries[rand3])
		}

		for _, e := range entries {
			e.Score = calScore(e)
			entries = append(entries, e)
		}

		sort.Sort(sort.Reverse(ByScore(entries)))

		for k := 1; k < (len(entries) - usedPopulationSize); k++ {
			entries = entries[:len(entries)-1] // Truncate slice.
		}
	}

	plotGraph(averageFitness, averageFitnessIndex)
}

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Less(i, j int) bool { return a[i].Score < a[j].Score }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func RandSeq(n int) string {
	var values = []rune("10")

	b := make([]rune, n)
	for i := range b {
		b[i] = values[rand.Intn(len(values))]
	}
	return string(b)
}

func calScore(e Entry) int {
	score := 0
	// // (i) & (iii)
	for k := 0; k < len(e.Value); k++ {
		if e.Value[k] == 49 { // eqaul to 1..... (rune)
			score++
		}
	}

	// // (ii)
	// targetE := Entry{Value: "01010001101000100011"}
	// for k := range e.Value {
	// 	if e.Value[k] == targetE.Value[k] {
	// 		score++
	// 	}
	// }

	// // (iii)
	if e.Value == "00000000000000000000" {
		score = 2 * maxFitness
	}

	return score
}

func mutation(e Entry) Entry {
	out := []rune(e.Value)
	i := rand.Intn(len(out))
	if out[i] == '1' {
		out[i] = '0'
	} else if out[i] == '0' {
		out[i] = '1'
	}

	e.Value = string(out)
	e.Score = calScore(e)
	return e
}

func crossover(e1, e2 Entry) (Entry, Entry) {

	rand := rand.Intn(len(e1.Value))
	temp1 := e1.Value[0:rand]
	tmep2 := e2.Value[0:rand]
	e1.Value = tmep2 + e1.Value[rand:]
	e2.Value = temp1 + e2.Value[rand:]
	e1.Score = calScore(e1)
	e2.Score = calScore(e2)

	return e1, e2
}

// Use to generate random binary strings
func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func plotGraph(results, resultIndex []float64) {

	p := plot.New()

	p.Title.Text = "Average Fitness per Generation"
	p.X.Label.Text = "Average Fitness"
	p.Y.Label.Text = "Generations"

	pts := make(plotter.XYs, len(results))

	for i := range results {

		pts[i].X = results[i]

		pts[i].Y = resultIndex[i]
	}

	plotutil.AddLinePoints(p, "", pts)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "Average_Fitness_per_Gen(iii).png"); err != nil {
		panic(err)
	}
}
