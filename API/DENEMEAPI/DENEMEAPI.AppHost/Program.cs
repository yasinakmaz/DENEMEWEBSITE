var builder = DistributedApplication.CreateBuilder(args);

builder.AddProject<Projects.DENEMEAPI>("denemeapi");

builder.Build().Run();
