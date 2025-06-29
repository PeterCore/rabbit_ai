import 'package:flutter/material.dart';

class SelectableBox extends StatelessWidget {
  final bool selected; // 是否选中
  final VoidCallback onTap; // 点击回调
  final Widget child; // 中间的内容，如图标或文字
  final Color selectedBorderColor;
  final Color unselectedBorderColor;
  final double borderWidth;
  final double borderRadius; // 圆角
  final EdgeInsets padding; // 内边距
  final Color? backgroundColor; // 背景色（可选）

  // 新增: tip 标签相关
  final String? tipLabel; // tip文字，如"推荐"
  final Color tipBackgroundColor; // tip背景色
  final Color tipTextColor; // tip文字颜色
  final EdgeInsets tipPadding; // tip内部边距
  final double tipBorderRadius; // tip圆角

  /// 使 tip 溢出到容器外面的偏移量，默认为 x = -8, y = -8
  final Offset tipOffset;

  const SelectableBox({
    super.key,
    required this.selected,
    required this.onTap,
    required this.child,
    this.selectedBorderColor = const Color(0xFFFF69B4),
    this.unselectedBorderColor = const Color(0xFFE0E0E0),
    this.borderWidth = 1.5,
    this.borderRadius = 8,
    this.padding = const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
    this.backgroundColor,
    // tip 属性
    this.tipLabel,
    this.tipBackgroundColor = Colors.red,
    this.tipTextColor = Colors.white,
    this.tipPadding = const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
    this.tipBorderRadius = 4,
    this.tipOffset = const Offset(4, -8),
  });

  @override
  Widget build(BuildContext context) {
    final borderColor = selected ? selectedBorderColor : unselectedBorderColor;

    return GestureDetector(
      onTap: onTap,
      child: Container(
        decoration: BoxDecoration(
          color: backgroundColor ?? Colors.white,
          border: Border.all(color: borderColor, width: borderWidth),
          borderRadius: BorderRadius.circular(borderRadius),
        ),
        // 使用 Stack 并允许溢出
        child: Stack(
          clipBehavior: Clip.none, // 允许子 Widget 溢出到父 Widget 区域外
          children: [
            // 主体内容
            Padding(
              padding: padding,
              child: child,
            ),

            // 如果有 tipLabel，就在右上角溢出
            if (tipLabel != null)
              Positioned(
                right: tipOffset.dx,
                top: tipOffset.dy,
                child: Container(
                  padding: tipPadding,
                  decoration: BoxDecoration(
                    color: tipBackgroundColor,
                    borderRadius: BorderRadius.circular(tipBorderRadius),
                  ),
                  child: Text(
                    tipLabel!,
                    style: TextStyle(color: tipTextColor, fontSize: 12),
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }
}
