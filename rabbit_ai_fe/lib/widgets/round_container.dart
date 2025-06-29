import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import '../app/app_color.dart';

// class RoundDashContainer extends StatelessWidget {
//   const RoundDashContainer(
//       {super.key,
//       required this.child,
//       this.backgroundColor,
//       this.horizontal,
//       this.color,
//       this.vertical});
//   final Widget child;
//   final Color? backgroundColor;
//   final double? horizontal;
//   final double? vertical;
//   final Color? color;
//   @override
//   Widget build(BuildContext context) {
//     return DottedBorder(
//         borderType: BorderType.RRect,
//         radius: Radius.circular(20.r),
//         dashPattern: const [3, 3], //
//         color: color ?? const Color(0xFF353F4A),
//         strokeWidth: 1.0,
//         child: Container(
//           decoration: BoxDecoration(
//             color: backgroundColor ?? const Color(0xFF353F4A),
//             borderRadius: BorderRadius.circular(20.r), // 圆角
//           ),
//           child: Padding(
//             padding: EdgeInsets.symmetric(
//                 horizontal: horizontal ?? 50.w, vertical: vertical ?? 25.h),
//             child: child,
//           ),
//         )); // 虚线宽度;
//   }
// }
//
// class RoundGradientContainer extends StatelessWidget {
//   const RoundGradientContainer(
//       {super.key,
//       this.child,
//       this.horizontal,
//       this.vertical,
//       this.boxShadow,
//       this.gradientColors,
//       this.backgroundColor,
//       this.borderColor,
//       this.borderWidth});
//   final Widget? child;
//   final double? horizontal;
//   final double? vertical;
//   final Color? backgroundColor;
//   final Color? borderColor;
//   final double? borderWidth;
//   final List<Color>? gradientColors;
//   final List<BoxShadow>? boxShadow;
//   @override
//   Widget build(BuildContext context) {
//     return Container(
//         padding: EdgeInsets.symmetric(
//             horizontal: horizontal ?? 24.w, vertical: vertical ?? 24.h),
//         decoration: BoxDecoration(
//           borderRadius: BorderRadius.circular(12.r),
//           border: Border.all(
//               strokeAlign: BorderSide.strokeAlignOutside,
//               style: BorderStyle.solid,
//               color: borderColor ?? Colors.white,
//               width: borderWidth ?? 2.w),
//           gradient: LinearGradient(
//             colors: gradientColors ?? [Colors.white, Colors.white],
//           ),
//         ),
//         child: child);
//   }
// }

class RoundContainer extends StatelessWidget {
  const RoundContainer({
    super.key,
    this.child,
    this.horizontal,
    this.vertical,
    this.boxShadow,
    this.imagePath,
    this.backgroundColor,
    this.radius,
    this.showBorder = false,
    this.showBorderColor,
    this.showBorderWidth,
    this.showBoxShadow = true,
    this.customBorderRadius, // ✅ 新增参数
  });

  final Widget? child;
  final double? horizontal;
  final double? vertical;
  final Color? backgroundColor;
  final String? imagePath;
  final double? radius;
  final bool showBorder;
  final bool? showBoxShadow;
  final Color? showBorderColor;
  final double? showBorderWidth;
  final List<BoxShadow>? boxShadow;

  /// ✅ 新增：可选的自定义圆角
  final BorderRadius? customBorderRadius;

  @override
  Widget build(BuildContext context) {
    final borderRadius =
        customBorderRadius ?? BorderRadius.circular(radius ?? 12.r);

    return imagePath == null
        ? Container(
            padding: EdgeInsets.symmetric(
                horizontal: horizontal ?? 24.w, vertical: vertical ?? 24.h),
            decoration: BoxDecoration(
              borderRadius: borderRadius,
              color: backgroundColor ?? AppColors.greyF8F8F8,
              border: (showBorder && showBorderColor != null)
                  ? Border.all(
                      color: showBorderColor ?? Colors.grey.withOpacity(.2),
                      width: showBorderWidth ?? 2,
                    )
                  : null,
              boxShadow: showBoxShadow == true
                  ? boxShadow ??
                      const [
                        BoxShadow(
                          color: Color(0xFFE6ECEA),
                          offset: Offset(0.0, 0.0),
                          blurRadius: 0,
                          spreadRadius: 0.5,
                        ),
                      ]
                  : null,
            ),
            child: child,
          )
        : Container(
            padding: EdgeInsets.symmetric(
                horizontal: horizontal ?? 24.w, vertical: vertical ?? 24.h),
            decoration: BoxDecoration(
              image: DecorationImage(
                image: AssetImage(imagePath ?? ""),
                fit: BoxFit.fill,
              ),
              borderRadius: borderRadius,
              color: backgroundColor ?? AppColors.greyF8F8F8,
              boxShadow: boxShadow ??
                  const [
                    BoxShadow(
                      color: Color(0xFFE6ECEA),
                      offset: Offset(0.0, 0.0),
                      blurRadius: 0,
                      spreadRadius: 0.5,
                    ),
                  ],
            ),
            child: child,
          );
  }
}
