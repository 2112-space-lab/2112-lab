from setuptools import setup, find_packages

setup(
    name="generator_graphql_api",
    version="0.0.2",
    description="GraphQL API Python Client",
    author="Elbujito",
    author_email="adrien.roques@icloud.com",
    url="https://github.com/2112-space-lab/2112-lab",
    packages=find_packages(exclude=["tests*"]),
    include_package_data=True,
    install_requires=[
        "ariadne",
        "ariadne-codegen",
        "fastapi",
        "uvicorn",
    ],
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
)
