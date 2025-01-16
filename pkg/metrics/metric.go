package metrics

import "github.com/prometheus/client_golang/prometheus"

// Создаем метрики
var (
	// // Время запроса каждого эндпоинта
	// HttpDurationSeconds = prometheus.NewHistogramVec(
	// 	prometheus.HistogramOpts{
	// 		Name:    "http_duration_seconds",
	// 		Help:    "Histogram of HTTP request durations by endpoint.",
	// 		Buckets: prometheus.DefBuckets,
	// 	},
	// 	[]string{"method", "endpoint"},
	// )

	// // Количество запросов каждого эндпоинта
	// HttpRequestsTotal = prometheus.NewCounterVec(
	// 	prometheus.CounterOpts{
	// 		Name: "http_requests_total",
	// 		Help: "Total number of HTTP requests by endpoint.",
	// 	},
	// 	[]string{"method", "endpoint"},
	// )

	// Время обращения в БД
	DbDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_duration_seconds",
			Help:    "Histogram of database access durations by method.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	// Время обращения во внешний API
	ApiDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_duration_seconds",
			Help:    "Histogram of external API access durations by method.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func MustAllRegister() {
	// Регистрируем метрики
	// prometheus.MustRegister(HttpDurationSeconds)
	// prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(DbDurationSeconds)
	prometheus.MustRegister(ApiDurationSeconds)
}
