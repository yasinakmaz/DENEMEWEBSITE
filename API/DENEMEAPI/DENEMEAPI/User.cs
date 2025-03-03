namespace DENEMEAPI
{
    public record User(
    int Id,
    string Username,
    string Email,
    string FullName,
    DateTime CreatedAt,
    bool IsActive
);
}
