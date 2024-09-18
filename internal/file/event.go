package file

// Event представляет событие, содержащее метрики и временную метку.
type Event struct {
	Gauge     map[string]float64 // Gauge хранит метрики типа gauge с их значениями.
	Counter   map[string]int64   // Counter хранит метрики типа counter с их значениями.
	Timestamp int64              // Timestamp содержит временную метку события.
}
