# kenetic-prometheus-exporter
Pretty simple Prometheus exporter for Keenetic router on golang.

Tested on `ZyXEL Keenetic Ultra II` with NDMS release `3.5.10`.


1) Based on web api polling(scrapping interval set to 10 seconds).
2) Collect `tx_load`, `rx_load`, `cpu_usage`, `free_ram`, `used_ram` metrics.
3) Export metrics with prometheus-exporter.

