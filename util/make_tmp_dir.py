import os

current_dir = os.getcwd()

with open(".env", "r") as f:
  lines = f.readlines()

for i, line in enumerate(lines):
  if line.startswith("TMP_DIR = "):
    lines[i] = f"TMP_DIR={os.path.join(current_dir, 'tmp')}\n"
    break

with open(".env", "w") as f:
  f.writelines(lines)

print(f"TMP_DIR is set to {os.path.join(current_dir, 'tmp')} in the .env file.")