namespace DENEMEAPI
{
    public class DashboardService : IDashboardService
    {
        private readonly ILogger<DashboardService> _logger;

        public DashboardService(ILogger<DashboardService> logger)
        {
            _logger = logger;
        }

        public async Task<DashboardData> GetDashboardDataAsync()
        {
            _logger.LogInformation("Dashboard verileri getiriliyor");

            // Gerçek uygulamada bu veriler veritabanından çekilirdi
            await Task.Delay(100); // Simüle edilmiş DB gecikmesi

            return new DashboardData(
                UsersTotal: 8942,
                ActiveUsers: 2531,
                NewToday: 147,
                Revenue: 134952.89m,
                RecentActivities: new List<Activity>
                {
                new Activity(1, "Ahmet Y.", "Yeni sipariş", "2 dakika önce", "₺1,299.00"),
                new Activity(2, "Mehmet K.", "Ödeme yapıldı", "15 dakika önce", "₺459.90"),
                new Activity(3, "Ayşe S.", "Hesap oluşturdu", "34 dakika önce", null),
                new Activity(4, "Fatma D.", "Sipariş iptali", "1 saat önce", "₺2,149.00"),
                new Activity(5, "Mustafa R.", "Yeni sipariş", "2 saat önce", "₺879.50")
                }
            );
        }
    }
}
