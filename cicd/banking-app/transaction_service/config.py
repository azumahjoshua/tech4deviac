import os
from pydantic_settings import BaseSettings, SettingsConfigDict
from pydantic import Field


class Settings(BaseSettings):
    DATABASE_URL: str = Field(..., alias="DATABASE_URL")
    CORE_API_URL: str = Field(default="http://localhost:8080")
    PORT: int = Field(default=5000)
    DEBUG: bool = Field(default=False)

    model_config = SettingsConfigDict(env_file=".env", extra="ignore")

    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        # Fix potential old-style postgres:// format
        if self.DATABASE_URL.startswith("postgres://"):
            self.DATABASE_URL = self.DATABASE_URL.replace("postgres://", "postgresql://", 1)


settings = Settings()
