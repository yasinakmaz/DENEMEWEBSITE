using System.Diagnostics;

namespace DENEMEAPI
{
    public record DashboardData(
    int UsersTotal,
    int ActiveUsers,
    int NewToday,
    decimal Revenue,
    IEnumerable<Activity> RecentActivities
);
}
