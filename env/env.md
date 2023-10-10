# Koa.env 介绍
一个`yaml`文件读取库

功能：
- 支持yaml配置文件读取
- 支持配置文件继承，合并子配置文件。
- 支持根据属性路径直接获取节点值

# 工作模式
```yaml
  // 测试文件 env.yaml
  test: 1
  test2: 
    test3: 1
    test4:
      - 1
      - 1
  // 测试文件 env.production.yaml
  test: 2
  test2: 
    test3: 2
    test4:
      - 2
      - 2
    test5: 2
  test3: aaaaa
```
```go
	path, _ := os.Getwd()
	err := env.LoadEnvConfig(path+"/env.yaml", "production")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(env.GetRootValue()) // => map[test:2 test2:map[test3:2 test4:[2 2] test5:2] test3:aaaaa]
	fmt.Println(env.GetValue("test2.test4[0]")) // => 2
```