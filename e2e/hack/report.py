import json
from pathlib import Path

with open('gateway-report.json') as f :
  report = json.load(f)

def generate_common_suite_output():
  passed_tests = []
  failed_tests = []
  for suite in report:
      for spec in suite['SpecReports']:
          if spec.get('LeafNodeType') != 'It':
              continue
          if spec['State'] == 'passed':
              passed_tests.append({
                  'test_name': spec['LeafNodeText'],
                  'location': f"{spec['LeafNodeLocation']['FileName']}:{spec['LeafNodeLocation']['LineNumber']}",
                  'start_time': spec['StartTime'],
                  'run_time': f"{spec['RunTime'] / 1e6:.2f}ms",  # 转换为毫秒
                  'events': [event['Message'] for event in spec['SpecEvents'] if event['SpecEventType'] == 'Node']
              })
          if spec['State'] == 'failed':
              failed_tests.append({
                  'test_name': spec['LeafNodeText'],
                  'location': f"{spec['LeafNodeLocation']['FileName']}:{spec['LeafNodeLocation']['LineNumber']}",
                  'start_time': spec['StartTime'],
                  'run_time': f"{spec['RunTime'] / 1e6:.2f}ms",  # 转换为毫秒
                  'events': [event['Message'] for event in spec['SpecEvents'] if event['SpecEventType'] == 'Node']
              })
  print("Gateway 测试结果")
  if passed_tests:
      print("  ✅ 成功用例列表")
      for test in passed_tests:
          print(f"  - 用例名称: {test['test_name']}")
          print(f"    代码位置: {test['location']}")
          print(f"    开始时间: {test['start_time']}")
          print(f"    运行时间: {test['run_time']}")
          print(f"    执行步骤: {', '.join(test['events'])}")

  if failed_tests:
      print("  ❌ 失败用例列表：")
      for test in failed_tests:
          print(f"  - 用例名称: {test['test_name']}")
          print(f"    代码位置: {test['location']}")
          print(f"    开始时间: {test['start_time']}")
          print(f"    运行时间: {test['run_time']}")
          print(f"    执行步骤: {', '.join(test['events'])}")
  print(f"成功用例{len(passed_tests)}\t失败用例{len(failed_tests)}")

def generate_test_coverage(output_file="TEST_COVERAGE.md"):
    """
    生成 Markdown 格式的测试覆盖文档

    :param output_file: 输出文件名
    """

    # 准备文档内容
    markdown = [
        "# 测试覆盖文档 (自动生成)",
        "",
        "| 模块 | 测试用例 | 状态 | 代码位置 |",
        "|------|---------|------|----------|"
    ]

    # 遍历所有测试套件
    for suite in report:
        for spec in suite.get('SpecReports', []):
            # 只处理具体的测试用例（忽略 BeforeSuite/AfterSuite 等）
            if spec.get('LeafNodeType') != 'It':
                continue

            # 提取模块名称（第一个层级文本）
            container_hierarchy = spec.get('ContainerHierarchyTexts', [])
            module = container_hierarchy[0] if container_hierarchy else "通用测试"

            # 提取用例信息
            test_case = {
                'name': spec.get('LeafNodeText', '未命名用例'),
                'status': '✅' if spec.get('State') == 'passed' else '❌',
                'file': spec.get('LeafNodeLocation', {}).get('FileName', ''),
                'line': spec.get('LeafNodeLocation', {}).get('LineNumber', '')
            }

            # 生成源码链接（支持 IDE 直接跳转）
            code_location = ""
            if test_case['file'] and test_case['line']:
                abs_path = Path(test_case['file']).absolute()
                code_location = f"[源码]({abs_path}#L{test_case['line']})"
            else:
                code_location = "`未记录位置`"

            # 生成表格行
            markdown.append(
                f"| {module} | {test_case['name']} | {test_case['status']} | {code_location} |"
            )

    # 写入文件
    try:
        with open(output_file, 'w', encoding='utf-8') as f:
            f.write('\n'.join(markdown))
        print(f"文档已生成：{Path(output_file).absolute()}")
    except Exception as e:
        print(f"写入文件失败: {str(e)}")

if __name__ == "__main__":
    # 使用示例
    generate_common_suite_output()
    generate_test_coverage(
        output_file="GATEWAY_TEST_COVERAGE.md"
    )
