import os
import re

directory = r"c:\Users\harsh\Python\StockTrack"

for root, dirs, files in os.walk(directory):
    for file in files:
        if file.endswith(".go"):
            filepath = os.path.join(root, file)
            
            with open(filepath, 'r', encoding='utf-8') as f:
                content = f.read()
            
            lines = content.split('\n')
            cleaned_lines = []
            
            for line in lines:
                stripped = line.lstrip()
                if stripped.startswith('//'):
                    continue
                
                if '//' in line and not ('"/' in line or '`/' in line):
                    idx = line.find('//')
                    line = line[:idx].rstrip()
                
                cleaned_lines.append(line)
            
            cleaned = '\n'.join(cleaned_lines)
            cleaned = re.sub(r'\n\n\n+', '\n\n', cleaned)
            
            with open(filepath, 'w', encoding='utf-8') as f:
                f.write(cleaned)
            
            print(f"Cleaned: {filepath}")

print("Done!")
