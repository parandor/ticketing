import argparse
import re
import os

class CodeExtractor:
    def __init__(self):
        self.type_pattern = r'type\s+(\w+)\s+struct\s*{([^}]*)}'
        self.empty_struct_pattern = r'type\s+(\w+)\s+struct\s*{}'
        self.interface_pattern = r'type\s+(\w+)\s+interface\s*{([^}]*)}'
        self.function_pattern = r'func\s+New(\w+)\(([^)]*)\)\s+(\w+)\s+\{([^}]+)\}'
        self.handler_pattern = r'func\s+New(\w+)Handler\(([^)]*)\)\s*\(([^)]*)\)\s*{([^}]*)}'

    def extract_code_blocks(self, input_text):
        type_matches = re.findall(self.type_pattern, input_text, re.DOTALL)
        empty_struct_matches = re.findall(self.empty_struct_pattern, input_text, re.DOTALL)
        interface_matches = re.findall(self.interface_pattern, input_text, re.DOTALL)
        function_matches = re.findall(self.function_pattern, input_text, re.DOTALL)
        handler_matches = re.findall(self.handler_pattern, input_text, re.DOTALL)

        extracted_code = []
        
        for match in type_matches:
            method_name = match[0]
            code_block = match[1].strip()
            extracted_code.append(("Type", f"type {method_name} struct {{\n{code_block}\n}}"))
            
        for match in empty_struct_matches:
            method_name = match[0]
            extracted_code.append(("Type", f"type {method_name} struct {{}}"))
            
        for match in interface_matches:
            method_name = match[0]
            code_block = match[1].strip()
            extracted_code.append(("Interface", f"type {method_name} interface {{\n{code_block}\n}}"))
            
        for match in function_matches:
            method_name = match[0]
            parameters = match[1]
            return_type = match[2]
            function_body = match[3]
            extracted_code.append(("Function", f"func New{method_name}({parameters}) {return_type} {{\n{function_body}\n}}"))
        
        for match in handler_matches:
            method_name = match[0]
            parameters1 = match[1]
            parameters2 = match[2]
            function_body = match[3]
            extracted_code.append(("Handler", f"func New{method_name}Handler({parameters1}) ({parameters2}) {{\n{function_body}\n}}"))
      
        return extracted_code

    def extract_code_blocks_from_file(self, file_path):
        with open(file_path, 'r') as file:
            input_text = file.read()
            return self.extract_code_blocks(input_text)

def main():
    parser = argparse.ArgumentParser(description='Extract code blocks from a file')
    parser.add_argument('file_path', type=str, help='Path to the file containing code blocks')
    args = parser.parse_args()

    extractor = CodeExtractor()
    file_path = os.path.join(os.getcwd(), args.file_path)
    extracted_code = extractor.extract_code_blocks_from_file(file_path)
    for _, code_block in extracted_code:
        print(code_block)
        print()

if __name__ == "__main__":
    main()
