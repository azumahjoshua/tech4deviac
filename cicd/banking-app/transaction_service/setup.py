from setuptools import setup, find_packages

setup(
    name="transaction_service",
    version="1.0",
    packages=find_packages(),
    install_requires=[
        'fastapi',
        'uvicorn',
        'sqlalchemy',
        'psycopg2-binary',
        'python-dotenv'
    ],
)