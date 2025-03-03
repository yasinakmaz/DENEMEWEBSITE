namespace DENEMEAPI
{
    public record Activity(
    int Id,
    string User,
    string Action,
    string Time,
    string? Amount
);
}
