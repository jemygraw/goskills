在运行任何 `goskills` 命令之前，您必须设置您的 OpenAI API 密钥：

```bash
export OPENAI_API_KEY="YOUR_OPENAI_API_KEY"
```

---

# goskills 使用示例

本文档提供了如何从各种编程语言执行 `goskills` 命令行工具的示例。

这些示例的基础命令是：

```shell
goskills run --auto-approve --model deepseek-v3 --api-base https://qianfan.baidubce.com/v2 --skills-dir=~/.claude/skills "使用markitdown 工具解析网页 https://baike.baidu.com/item/%E5%AD%94%E5%AD%90/1584"
```

---

## Shell (Bash)

这是运行该命令最直接的方式。

```bash
#!/bin/bash

# 定义提示
PROMPT="使用markitdown 工具解析网页 https://baike.baidu.com/item/%E5%AD%94%E5%AD%90/1584"

# 执行命令并捕获输出
RESULT=$(goskills run --auto-approve --model deepseek-v3 --api-base https://qianfan.baidubce.com/v2 --skills-dir=$HOME/.claude/skills "$PROMPT")

# 或者直接执行并打印
# goskills run --auto-approve --model deepseek-v3 --api-base https://qianfan.baidubce.com/v2 --skills-dir=$HOME/.claude/skills "$PROMPT"


echo "输出:"
echo "$RESULT"
```

---

## Python

使用 `subprocess` 模块是在 Python 中运行外部命令的标准方法。为了正确处理 `~` 符号，我们需要使用 `os.path.expanduser` 来获取完整路径。

```python
import subprocess
import shlex
import os # 导入 os 模块

# 定义提示
prompt = "使用markitdown 工具解析网页 https://baike.baidu.com/item/%E5%AD%94%E5%AD%90/1584"

# 获取并扩展技能目录路径
skills_dir_path = os.path.expanduser("~/.claude/skills")
skills_dir_arg = f"--skills-dir={skills_dir_path}"

# 为安全起见，将命令定义为参数列表
command = [
    "goskills", "run",
    "--auto-approve",
    "--model", "deepseek-v3",
    "--api-base", "https://qianfan.baidubce.com/v2",
    skills_dir_arg, # 使用构建的路径参数
    prompt
]

# 或者，使用 shlex 从字符串构建命令以进行正确的引用
# 注意：如果使用 shlex.split，需要确保 skills_dir_path 已正确扩展，且整个字符串被正确引用
# 例如：
# cmd_str = f'goskills run --auto-approve --model deepseek-v3 --api-base https://qianfan.baidubce.com/v2 {skills_dir_arg} "{prompt}"'
# command = shlex.split(cmd_str)


try:
    # 执行命令，捕获 stdout 和 stderr
    result = subprocess.run(
        command,
        check=True,        # 如果命令返回非零退出码，则引发异常
        capture_output=True, # 捕获 stdout 和 stderr
        text=True          # 将 stdout/stderr 解码为文本
    )
    
    print("命令执行成功:")
    print("输出:\n", result.stdout)

except FileNotFoundError:
    print("错误: 'goskills' 命令未找到。请确保它在您的 PATH 中。")
except subprocess.CalledProcessError as e:
    print(f"命令失败，退出码 {e.returncode}:")
    print("标准错误输出:\n", e.stderr)

```

---

## JavaScript (Node.js)

在 Node.js 中，您可以使用 `child_process` 模块。

```javascript
const { exec } = require('child_process');

// 使用单引号作为外部字符串以便轻松处理内部的双引号
const command = 'goskills run --auto-approve --model deepseek-v3 --api-base https://qianfan.baidubce.com/v2 --skills-dir=$HOME/.claude/skills "使用markitdown 工具解析网页 https://baike.baidu.com/item/%E5%AD%94%E5%AD%90/1584"';

exec(command, (error, stdout, stderr) => {
    if (error) {
        console.error(`执行错误: ${error.message}`);
        if (stderr) {
            console.error(`标准错误输出: ${stderr}`);
        }
        return;
    }

    console.log(`命令输出:\n${stdout}`);
});
```

---

## Go

在 Go 中，使用 `os/exec` 包来运行外部命令。为了正确处理 `~` 符号，我们需要手动获取用户主目录并构建路径。

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	prompt := "使用markitdown 工具解析网页 https://baike.baidu.com/item/%E5%AD%94%E5%AD%90/1584"

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("无法获取主目录: %v\n", err)
		return
	}
	// 构建技能目录的完整路径
	skillsDirPath := filepath.Join(homeDir, ".claude", "skills")
	skillsDirArg := "--skills-dir=" + skillsDirPath
	
	cmd := exec.Command("goskills", "run",
		"--auto-approve",
		"--model", "deepseek-v3",
		"--api-base", "https://qianfan.baidubce.com/v2",
		skillsDirArg, // 使用构建的路径参数
		prompt)

	// CombinedOutput 运行命令并返回其组合的标准输出和标准错误
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("执行命令时出错: %v\n", err)
		fmt.Printf("输出:\n%s\n", string(output))
		return
	}

	fmt.Printf("命令输出:\n%s\n", string(output))
}
```

---

## Java

使用 `ProcessBuilder` 是在 Java 中执行命令的现代和推荐方式。为了正确指定 `--skills-dir` 参数，我们需要手动获取用户的主目录并构建完整路径。

```java
import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.IOException;
import java.nio.file.Paths; // 导入 Paths 类

public class GoSkillsRunner {
    public static void main(String[] args) {
        try {
            // 获取用户主目录并构建技能目录路径
            String userHome = System.getProperty("user.home");
            String skillsDirPath = Paths.get(userHome, ".claude", "skills").toString();
            String skillsDirArg = "--skills-dir=" + skillsDirPath;

            ProcessBuilder pb = new ProcessBuilder(
                "goskills", "run",
                "--auto-approve",
                "--model", "deepseek-v3",
                "--api-base", "https://qianfan.baidubce.com/v2",
                skillsDirArg, // 使用构建的路径参数
                "使用markitdown 工具解析网页 https://baike.baidu.com/item/%E5%AD%94%E5%AD%90/1584"
            );

            // 将错误流重定向到与标准输出流相同
            pb.redirectErrorStream(true);

            Process process = pb.start();

            // 从命令读取输出
            StringBuilder output = new StringBuilder();
            try (BufferedReader reader = new BufferedReader(new InputStreamReader(process.getInputStream()))) {
                String line;
                while ((line = reader.readLine()) != null) {
                    output.append(line).append("\n");
                }
            }

            int exitCode = process.waitFor();
            System.out.println("退出码: " + exitCode);
            System.out.println("输出:\n" + output.toString());

        } catch (IOException | InterruptedException e) {
            e.printStackTrace();
        }
    }
}
```

---

## Rust

在 Rust 中，使用 `std::process::Command` 结构体来执行外部命令。为了正确处理 `$HOME` 环境变量，我们需要手动获取它的值并构建路径。

```rust
use std::process::Command;
use std::env;
use std::path::PathBuf;

fn main() {
    let prompt = "使用markitdown 工具解析网页 https://baike.baidu.com/item/%E5%AD%94%E5%AD%90/1584";

    // 获取 HOME 环境变量
    let home_dir = env::var("HOME").expect("HOME 环境变量未设置");
    // 构建技能目录路径
    let skills_dir_path = PathBuf::from(home_dir).join(".claude").join("skills");
    let skills_dir_arg = format!("--skills-dir={}", skills_dir_path.to_str().expect("路径无效"));

    let output = Command::new("goskills")
        .arg("run")
        .arg("--auto-approve")
        .arg("--model")
        .arg("deepseek-v3")
        .arg("--api-base")
        .arg("https://qianfan.baidubce.com/v2")
        .arg(&skills_dir_arg) // 使用构建的路径参数
        .arg(prompt)
        .output()
        .expect("未能执行命令");

    if output.status.success() {
        println!("命令执行成功:");
        println!("输出:\n{}", String::from_utf8_lossy(&output.stdout));
    } else {
        eprintln!("命令失败，退出码: {:?}", output.status.code());
        eprintln!("标准错误输出:\n{}", String::from_utf8_lossy(&output.stderr));
    }
}
```

---

## C++

在 C++ 中，`std::system` 提供了一种简单的方式来执行 shell 命令。然而，为了捕获命令的输出，`popen` 是一个更好的选择。`popen` 会执行一个命令并创建一个管道，允许程序读取该命令的标准输出。

```cpp
#include <iostream>
#include <cstdio>   // For popen, pclose, fgets
#include <string>
#include <array>

int main() {
    std::string prompt = "使用markitdown 工具解析网页 https://baike.baidu.com/item/%E5%AD%94%E5%AD%90/1584";
    std::string command = "goskills run --auto-approve --model deepseek-v3 --api-base https://qianfan.baidubce.com/v2 --skills-dir=$HOME/.claude/skills \"" + prompt + "\"";

    // 使用 "r" 模式执行 popen 以读取命令的输出
    FILE* pipe = popen(command.c_str(), "r");
    if (!pipe) {
        std::cerr << "无法执行命令！" << std::endl;
        return 1;
    }

    std::array<char, 256> buffer;
    std::string result;
    while (fgets(buffer.data(), buffer.size(), pipe) != nullptr) {
        result += buffer.data();
    }

    // pclose 会等待命令终止并返回其退出状态
    int exit_code = pclose(pipe);
    
    std::cout << "命令输出:" << std::endl;
    std::cout << result << std::endl;
    std::cout << "退出码: " << WEXITSTATUS(exit_code) << std::endl;


    return 0;
}
```

---

## C

在 C 语言中，`system()` 函数是 `stdlib.h` 的一部分，允许执行 shell 命令。但它不容易捕获输出。为了读取命令的输出，推荐使用 `popen` 函数，它会创建一个管道来连接到被调用进程的标准输出。

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define BUFFER_SIZE 256

int main() {
    char prompt[] = "使用markitdown 工具解析网页 https://baike.baidu.com/item/%E5%AD%94%E5%AD%90/1584";
    char command[512];
    snprintf(command, sizeof(command), "goskills run --auto-approve --model deepseek-v3 --api-base https://qianfan.baidubce.com/v2 --skills-dir=$HOME/.claude/skills \"%s\"", prompt);

    FILE *pipe = popen(command, "r");
    if (pipe == NULL) {
        fprintf(stderr, "无法执行命令！\n");
        return 1;
    }

    char buffer[BUFFER_SIZE];
    printf("命令输出:\n");
    // 逐行读取管道的输出并打印
    while (fgets(buffer, sizeof(buffer), pipe) != NULL) {
        printf("%s", buffer);
    }

    // pclose 等待命令终止并返回其退出状态
    int exit_code = pclose(pipe);
    fprintf(stdout, "\n退出码: %d\n", WEXITSTATUS(exit_code));

    return 0;
}
```