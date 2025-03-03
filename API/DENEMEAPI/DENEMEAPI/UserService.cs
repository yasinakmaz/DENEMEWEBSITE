namespace DENEMEAPI
{
    public class UserService : IUserService
    {
        private readonly List<User> _users;
        private readonly ILogger<UserService> _logger;
        private int _nextId = 6;

        public UserService(ILogger<UserService> logger)
        {
            _logger = logger;
            _users = new List<User>
        {
            new User(1, "ahmet123", "ahmet@ornek.com", "Ahmet Yılmaz", DateTime.Now.AddDays(-45), true),
            new User(2, "mehmetk", "mehmet@ornek.com", "Mehmet Kaya", DateTime.Now.AddDays(-30), true),
            new User(3, "ayses", "ayse@ornek.com", "Ayşe Şahin", DateTime.Now.AddDays(-10), true),
            new User(4, "fatmad", "fatma@ornek.com", "Fatma Demir", DateTime.Now.AddDays(-5), true),
            new User(5, "mustafar", "mustafa@ornek.com", "Mustafa Reis", DateTime.Now.AddDays(-2), true)
        };
        }

        public async Task<IEnumerable<User>> GetAllUsersAsync()
        {
            _logger.LogInformation("Tüm kullanıcılar getiriliyor");
            await Task.Delay(50); // Simüle edilmiş gecikme
            return _users;
        }

        public async Task<User?> GetUserByIdAsync(int id)
        {
            _logger.LogInformation("Kullanıcı ID ile getiriliyor: {Id}", id);
            await Task.Delay(30); // Simüle edilmiş gecikme
            return _users.Find(u => u.Id == id);
        }

        public async Task<User> CreateUserAsync(User user)
        {
            _logger.LogInformation("Yeni kullanıcı oluşturuluyor: {Username}", user.Username);

            // Gerçek uygulamada doğrulama kontrolleri yapılır
            var newUser = user with { Id = _nextId++, CreatedAt = DateTime.Now };
            _users.Add(newUser);

            await Task.Delay(100); // Simüle edilmiş gecikme
            return newUser;
        }
    }
}
