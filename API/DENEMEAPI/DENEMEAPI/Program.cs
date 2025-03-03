using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using System.Text.Json;
using System.Text.Json.Serialization;
using System;
using System.Collections.Generic;
using System.Threading.Tasks;

var builder = WebApplication.CreateBuilder(args);

// CORS Politikas� Ekleme
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowAll", builder =>
    {
        builder.AllowAnyOrigin()
               .AllowAnyMethod()
               .AllowAnyHeader();
    });
});

// JSON Se�enekleri
builder.Services.ConfigureHttpJsonOptions(options =>
{
    options.SerializerOptions.PropertyNamingPolicy = JsonNamingPolicy.CamelCase;
    options.SerializerOptions.DefaultIgnoreCondition = JsonIgnoreCondition.WhenWritingNull;
});

// Swagger API Dok�mantasyonu
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

// Dependency Injection i�in Servisleri Ekleme
builder.Services.AddSingleton<IDashboardService, DashboardService>();
builder.Services.AddSingleton<IUserService, UserService>();

var app = builder.Build();

// Middleware Yap�land�rmas�
app.UseSwagger();
app.UseSwaggerUI();
app.UseCors("AllowAll");

// API Rotalar� Tan�mlama
app.MapGet("/api/dashboard", async (IDashboardService dashboardService) =>
{
    var dashboardData = await dashboardService.GetDashboardDataAsync();
    return Results.Ok(dashboardData);
})
.WithName("GetDashboardData")
.WithOpenApi();

app.MapGet("/api/users", async (IUserService userService) =>
{
    var users = await userService.GetAllUsersAsync();
    return Results.Ok(users);
})
.WithName("GetAllUsers")
.WithOpenApi();

app.MapGet("/api/users/{id}", async (int id, IUserService userService) =>
{
    var user = await userService.GetUserByIdAsync(id);
    if (user == null)
    {
        return Results.NotFound($"Kullan�c� bulunamad� (ID: {id})");
    }
    return Results.Ok(user);
})
.WithName("GetUserById")
.WithOpenApi();

app.MapPost("/api/users", async (User user, IUserService userService) =>
{
    var newUser = await userService.CreateUserAsync(user);
    return Results.Created($"/api/users/{newUser.Id}", newUser);
})
.WithName("CreateUser")
.WithOpenApi();

// .NET Aspire i�in telemetri konfig�rasyonu
if (builder.Environment.IsDevelopment())
{
    app.MapPrometheusScrapingEndpoint();
    app.UseAspireInstrumentation();
}

// API'yi Ba�lat
app.Run();