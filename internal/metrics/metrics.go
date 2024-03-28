package metrics

import (
	"sync"

	"github.com/henrywhitaker3/adguard-exporter/internal/adguard"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ScrapeErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "scrape_errors_total",
		Namespace: "adguard",
		Help:      "The number of errors scraping a target",
	}, []string{"server"})

	// Status
	ProtectionEnabled = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "protection_enabled",
		Namespace: "adguard",
		Help:      "Whether DNS filtering is enabled",
	}, []string{"server"})
	Running = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "running",
		Namespace: "adguard",
		Help:      "Whether adguard is running or not",
	}, []string{"server", "version"})

	// Stats
	TotalQueries = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "queries",
		Namespace: "adguard",
		Help:      "Total queries processed in the last 24 hours",
	}, []string{"server"})
	BlockedFiltered = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "queries_blocked",
		Namespace: "adguard",
		Help:      "Total queries that have been blocked from filter lists",
	}, []string{"server"})
	ReplacedSafesearch = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "replaced_safesearch",
		Namespace: "adguard",
		Help:      "Total queries that have been replaced due to safesearch",
	}, []string{"server"})
	ReplacedSafebrowsing = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "replaced_safebrowsing",
		Namespace: "adguard",
		Help:      "Total queries that have been replaced due to safebrowsing",
	}, []string{"server"})
	ReplcaedParental = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "replaced_parental",
		Namespace: "adguard",
		Help:      "Total queries that have been replaced due to parental",
	}, []string{"server"})
	AvgProcessingTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "avg_processing_time_seconds",
		Namespace: "adguard",
		Help:      "The average query processing time in seconds",
	}, []string{"server"})
	TopQueriedDomains = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "top_queried_domains",
		Namespace: "adguard",
		Help:      "The number of queries for the top domains",
	}, []string{"server", "domain"})
	TopBlockedDomains = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "top_blocked_domains",
		Namespace: "adguard",
		Help:      "The number of blocked queries for the top domains",
	}, []string{"server", "domain"})
	TopClients = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "top_clients",
		Namespace: "adguard",
		Help:      "The number of queries for the top clients",
	}, []string{"server", "client"})
	TopUpstreams = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "top_upstreams",
		Namespace: "adguard",
		Help:      "The number of repsonses for the top upstream servers",
	}, []string{"server", "upstream"})
	TopUpstreamsAvgTimes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "top_upstreams_avg_response_time_seconds",
		Namespace: "adguard",
		Help:      "The average response time for each of the top upstream servers",
	}, []string{"server", "upstream"})
	QueryTypes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "query_types",
		Namespace: "adguard",
		Help:      "The number of queries for a specific type",
	}, []string{"server", "type"})

	// DHCP
	DhcpEnabled = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "dhcp_enabled",
		Namespace: "adguard",
		Help:      "Whether dhcp is enabled",
	}, []string{"server"})
	DhcpLeasesMetric = prometheus.NewDesc(
		"adguard_dhcp_leases",
		"The dhcp leases",
		[]string{"server", "type", "ip", "mac", "hostname"},
		nil,
	)
	DhcpLeases = NewDhcpLeasesServer(DhcpLeasesMetric)
)

type DhcpLeasesServer struct {
	mu     *sync.Mutex
	Desc   *prometheus.Desc
	leases map[string][]adguard.DhcpLease
}

func NewDhcpLeasesServer(desc *prometheus.Desc) *DhcpLeasesServer {
	return &DhcpLeasesServer{
		mu:     &sync.Mutex{},
		leases: map[string][]adguard.DhcpLease{},
		Desc:   desc,
	}
}

func (d *DhcpLeasesServer) Record(server string, leases []adguard.DhcpLease) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.leases[server] = leases
}

func (d *DhcpLeasesServer) Collect(ch chan<- prometheus.Metric) {
	for server, leases := range d.leases {
		for _, lease := range leases {
			ch <- prometheus.MustNewConstMetric(
				d.Desc,
				prometheus.CounterValue,
				1,
				server, lease.Type, lease.IP, lease.Mac, lease.Hostname,
			)
		}
	}
}

func (d *DhcpLeasesServer) Describe(ch chan<- *prometheus.Desc) {
	ch <- d.Desc
}

func Init() {
	prometheus.MustRegister(ScrapeErrors)

	// Stats
	prometheus.MustRegister(TotalQueries)
	prometheus.MustRegister(BlockedFiltered)
	prometheus.MustRegister(ReplacedSafesearch)
	prometheus.MustRegister(ReplacedSafebrowsing)
	prometheus.MustRegister(ReplcaedParental)
	prometheus.MustRegister(AvgProcessingTime)
	prometheus.MustRegister(TopBlockedDomains)
	prometheus.MustRegister(TopClients)
	prometheus.MustRegister(TopQueriedDomains)
	prometheus.MustRegister(TopUpstreams)
	prometheus.MustRegister(TopUpstreamsAvgTimes)
	prometheus.MustRegister(QueryTypes)

	// Status
	prometheus.MustRegister(Running)
	prometheus.MustRegister(ProtectionEnabled)

	// DHCP
	prometheus.MustRegister(DhcpEnabled)
	prometheus.MustRegister(DhcpLeases)
}
