namespace DENEMEAPI
{
    public interface IDashboardService
    {
        Task<DashboardData> GetDashboardDataAsync();
    }
}
