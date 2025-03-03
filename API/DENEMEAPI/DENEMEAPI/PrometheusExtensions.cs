namespace DENEMEAPI
{
    public static class PrometheusExtensions
    {
        public static void MapPrometheusScrapingEndpoint(this WebApplication app)
        {
            // Prometheus metrik toplamak için endpoint
            app.MapGet("/metrics", () => "Prometheus metrics would be here");
        }
    }
}
