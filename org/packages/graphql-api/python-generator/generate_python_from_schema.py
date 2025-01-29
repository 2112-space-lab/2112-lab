import os
import glob
from graphql import build_schema
from graphql.utilities import print_schema

SCHEMA_PATH = "./org/graphql-api/schemas/*.graphqls"  # Ensure this path matches your schema directory
OUTPUT_PATH = "./org/graphql-api/python-generator/generated/schema.py"

def generate_python_code():
    # Ensure the output directory exists
    output_dir = os.path.dirname(OUTPUT_PATH)
    os.makedirs(output_dir, exist_ok=True)

    # Collect all schema files
    schema_files = glob.glob(SCHEMA_PATH)
    if not schema_files:
        raise FileNotFoundError(f"No GraphQL schema files found in {SCHEMA_PATH}")

    # Concatenate the contents of all schema files
    schema_content = ""
    for file_path in schema_files:
        with open(file_path, "r") as schema_file:
            content = schema_file.read().strip()
            if content:
                schema_content += content + "\n"

    # Check if the schema content is valid
    if not schema_content.strip():
        raise ValueError("Schema content is empty. Ensure schema files are valid and not empty.")

    # Build the schema and generate Python code
    schema = build_schema(schema_content)
    python_code = f'"""Generated GraphQL Schema"""\n\n{print_schema(schema)}'

    # Write the generated code to the output file
    with open(OUTPUT_PATH, "w") as output_file:
        output_file.write(python_code)

    print(f"Python GraphQL schema generated at {OUTPUT_PATH}")

if __name__ == "__main__":
    try:
        generate_python_code()
    except Exception as e:
        print(f"Error: {e}")
        raise
