package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bytebotgo/gcoord-go/gcoord"
	"github.com/spf13/cobra"
)

// 文本格式化函数（为保持向后兼容性保留）
func red(s string) string     { return s }
func green(s string) string   { return s }
func blue(s string) string    { return s }
func yellow(s string) string  { return s }
func cyan(s string) string    { return s }
func magenta(s string) string { return s }
func bold(s string) string    { return s }

// 版本信息
const version = "1.0.0"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gcoord",
		Short: "地理坐标转换工具",
		Long: fmt.Sprintf(`%s gcoord-go 地理坐标转换工具

一个功能强大的地理坐标转换命令行工具，支持多种坐标系之间的转换。

支持的坐标系:
  %s  WGS84: 世界大地坐标系，GPS原始坐标
  %s  GCJ02: 国测局坐标系，中国标准坐标系  
  %s  BD09: 百度坐标系
  %s  BD09MC: 百度墨卡托投影坐标系
  %s  EPSG3857: Web墨卡托投影坐标系

使用示例:
  %s 转换单个坐标点
  %s 转换JSON格式的坐标
  %s 显示支持的坐标系

更多信息请使用子命令的 --help 参数查看。`,
			bold("🗺️"),
			yellow("•"),
			yellow("•"),
			yellow("•"),
			yellow("•"),
			yellow("•"),
			green("gcoord convert -from WGS84 -to GCJ02 -lon 116.397 -lat 39.908"),
			green(`gcoord convert -from WGS84 -to BD09 -json '{"type":"Point","coordinates":[116.397,39.908]}'`),
			green("gcoord list"),
		),
		Version: version,
	}

	// 添加子命令
	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(listCmd)

	// 设置版本信息
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s gcoord-go v{{.Version}}\n", bold("🗺️")))

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%s 执行错误: %v\n", red("❌"), err)
		os.Exit(1)
	}
}

// convertCmd 坐标转换命令
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "转换地理坐标",
	Long: fmt.Sprintf(`%s 坐标转换命令

将坐标从一个坐标系转换到另一个坐标系。

支持的输入格式:
  %s 经纬度坐标: -lon <经度> -lat <纬度>
  %s JSON格式: -json '<JSON坐标>'
  %s GeoJSON对象: Point, LineString, Polygon, Feature, FeatureCollection

示例:
  %s 转换单个坐标点
  %s 转换JSON格式的坐标
  %s 转换GeoJSON Feature`,
		bold("🔄"),
		yellow("•"),
		yellow("•"),
		yellow("•"),
		green("gcoord convert -from WGS84 -to GCJ02 -lon 116.397 -lat 39.908"),
		green(`gcoord convert -from WGS84 -to BD09 -json '{"type":"Point","coordinates":[116.397,39.908]}'`),
		green(`gcoord convert -from GCJ02 -to EPSG3857 -json '{"type":"Feature","geometry":{"type":"Point","coordinates":[116.397,39.908]}}'`),
	),
	Run: runConvert,
}

// listCmd 列出支持的坐标系
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "显示支持的坐标系",
	Long: fmt.Sprintf(`%s 坐标系列表

显示所有支持的坐标系及其详细信息。`,
		bold("📋"),
	),
	Run: runList,
}

func init() {
	// convert 命令参数
	convertCmd.Flags().StringP("from", "f", "", "源坐标系 (必需)")
	convertCmd.Flags().StringP("to", "t", "", "目标坐标系 (必需)")
	convertCmd.Flags().Float64("lon", 0, "经度")
	convertCmd.Flags().Float64("lat", 0, "纬度")
	convertCmd.Flags().StringP("json", "j", "", "JSON格式的坐标输入")
	convertCmd.Flags().BoolP("verbose", "v", false, "显示详细信息")

	// 标记必需参数
	convertCmd.MarkFlagRequired("from")
	convertCmd.MarkFlagRequired("to")
}

func runConvert(cmd *cobra.Command, args []string) {
	// 获取参数
	fromCRS, _ := cmd.Flags().GetString("from")
	toCRS, _ := cmd.Flags().GetString("to")
	lon, _ := cmd.Flags().GetFloat64("lon")
	lat, _ := cmd.Flags().GetFloat64("lat")
	jsonInput, _ := cmd.Flags().GetString("json")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// 验证坐标系
	if !isValidCRS(fromCRS) {
		fmt.Printf("%s 错误: 无效的源坐标系 '%s'\n", red("❌"), fromCRS)
		showValidCRS()
		os.Exit(1)
	}

	if !isValidCRS(toCRS) {
		fmt.Printf("%s 错误: 无效的目标坐标系 '%s'\n", red("❌"), toCRS)
		showValidCRS()
		os.Exit(1)
	}

	// 处理输入
	var input interface{}
	var err error

	if jsonInput != "" {
		// JSON 输入
		input, err = parseJSONInput(jsonInput)
		if err != nil {
			fmt.Printf("%s JSON解析错误: %v\n", red("❌"), err)
			os.Exit(1)
		}
	} else {
		// 经纬度输入
		if lon == 0 && lat == 0 {
			fmt.Printf("%s 错误: 必须提供经纬度坐标或JSON输入\n", red("❌"))
			fmt.Printf("使用 %s 查看帮助信息\n", cyan("gcoord convert --help"))
			os.Exit(1)
		}
		input = gcoord.Position{lon, lat}
	}

	// 执行转换
	result, err := gcoord.Transform(input, gcoord.CRSTypes(fromCRS), gcoord.CRSTypes(toCRS))
	if err != nil {
		fmt.Printf("%s 转换错误: %v\n", red("❌"), err)
		os.Exit(1)
	}

	// 显示结果
	showResult(input, result, fromCRS, toCRS, verbose)
}

func runList(cmd *cobra.Command, args []string) {
	fmt.Printf("%s 支持的坐标系\n\n", bold("📋"))

	crsList := []struct {
		name        string
		description string
		aliases     []string
		precision   string
	}{
		{
			name:        "WGS84",
			description: "世界大地坐标系，GPS原始坐标",
			aliases:     []string{"WGS1984", "EPSG4326"},
			precision:   "约 1e-5 度 (约 1 米)",
		},
		{
			name:        "GCJ02",
			description: "国测局坐标系，中国标准坐标系",
			aliases:     []string{"AMap"},
			precision:   "约 1e-5 度 (约 1 米)",
		},
		{
			name:        "BD09",
			description: "百度坐标系",
			aliases:     []string{"BD09LL", "Baidu", "BMap"},
			precision:   "约 1e-5 度 (约 1 米)",
		},
		{
			name:        "BD09MC",
			description: "百度墨卡托投影坐标系",
			aliases:     []string{"BD09Meter"},
			precision:   "约 1 米",
		},
		{
			name:        "EPSG3857",
			description: "Web墨卡托投影坐标系",
			aliases:     []string{"EPSG900913", "EPSG102100", "WebMercator", "WM"},
			precision:   "约 1 米",
		},
	}

	for _, crs := range crsList {
		fmt.Printf("%s %s\n", blue("📍"), bold(crs.name))
		fmt.Printf("  描述: %s\n", crs.description)
		if len(crs.aliases) > 0 {
			fmt.Printf("  别名: %s\n", strings.Join(crs.aliases, ", "))
		}
		fmt.Printf("  精度: %s\n", cyan(crs.precision))
		fmt.Println()
	}

	fmt.Printf("%s 转换路径支持:\n", bold("🔄"))
	fmt.Printf("  %s 所有坐标系之间都可以相互转换\n", green("✓"))
	fmt.Printf("  %s 自动选择最优转换路径\n", green("✓"))
	fmt.Printf("  %s 支持链式转换 (如: WGS84 → GCJ02 → BD09)\n", green("✓"))
}

func isValidCRS(crs string) bool {
	validCRS := map[string]bool{
		"WGS84":    true,
		"GCJ02":    true,
		"BD09":     true,
		"BD09MC":   true,
		"EPSG3857": true,
	}
	return validCRS[crs]
}

func parseJSONInput(jsonStr string) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	return result, err
}

func showResult(input, result interface{}, fromCRS, toCRS string, verbose bool) {
	if verbose {
		fmt.Printf("\n%s 坐标转换结果\n", bold("🎯"))
		fmt.Printf("%s\n", strings.Repeat("=", 50))

		// 显示输入信息
		fmt.Printf("\n%s 源坐标系: %s\n", blue("📍"), magenta(fromCRS))
		fmt.Printf("%s 目标坐标系: %s\n", blue("🎯"), magenta(toCRS))

		// 显示输入坐标
		fmt.Printf("\n%s 输入坐标:\n", yellow("📥"))
		showCoordinate(input, "  ")

		// 显示输出坐标
		fmt.Printf("\n%s 输出坐标:\n", green("📤"))
		showCoordinate(result, "  ")
	} else {
		// 简洁输出
		if pos, ok := result.(gcoord.Position); ok && len(pos) >= 2 {
			fmt.Printf("%.6f,%.6f\n", pos[0], pos[1])
			return
		}
	}

	// 如果是简单坐标点，显示格式化输出
	if pos, ok := result.(gcoord.Position); ok && len(pos) >= 2 {
		if verbose {
			fmt.Printf("\n%s 格式化输出:\n", cyan("📋"))
			fmt.Printf("  经度: %s\n", green(fmt.Sprintf("%.6f", pos[0])))
			fmt.Printf("  纬度: %s\n", green(fmt.Sprintf("%.6f", pos[1])))

			// 显示精度信息
			fmt.Printf("\n%s 转换精度:\n", blue("🎯"))
			if toCRS == "EPSG3857" || toCRS == "BD09MC" {
				fmt.Printf("  投影坐标精度: 约 1 米\n")
			} else {
				fmt.Printf("  经纬度精度: 约 1e-5 度 (约 1 米)\n")
			}

			fmt.Printf("\n%s\n", strings.Repeat("=", 50))
		}
	}
}

func showCoordinate(coord interface{}, prefix string) {
	switch c := coord.(type) {
	case gcoord.Position:
		if len(c) >= 2 {
			fmt.Printf("%s[%.6f, %.6f]\n", prefix, c[0], c[1])
		} else {
			fmt.Printf("%s%v\n", prefix, c)
		}
	case []float64:
		if len(c) >= 2 {
			fmt.Printf("%s[%.6f, %.6f]\n", prefix, c[0], c[1])
		} else {
			fmt.Printf("%s%v\n", prefix, c)
		}
	case string:
		// 尝试格式化JSON
		var obj interface{}
		if err := json.Unmarshal([]byte(c), &obj); err == nil {
			formatted, _ := json.MarshalIndent(obj, prefix, "  ")
			fmt.Printf("%s\n", string(formatted))
		} else {
			fmt.Printf("%s%s\n", prefix, c)
		}
	default:
		// 尝试转换为JSON
		if jsonBytes, err := json.MarshalIndent(c, prefix, "  "); err == nil {
			fmt.Printf("%s\n", string(jsonBytes))
		} else {
			fmt.Printf("%s%v\n", prefix, c)
		}
	}
}

func showValidCRS() {
	validCRS := []string{"WGS84", "GCJ02", "BD09", "BD09MC", "EPSG3857"}
	fmt.Printf("支持的坐标系: %s\n", strings.Join(validCRS, ", "))
}
