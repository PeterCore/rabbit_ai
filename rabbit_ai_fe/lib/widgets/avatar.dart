import 'package:flutter/material.dart';

import '../app/app_color.dart';
import '../app/app_style.dart';
import 'net_image.dart';

class Avatar extends StatelessWidget {
  final String? url;
  final bool showBorder;
  final double size;
  final double? radius;
  final Color? showBorderColor;
  final double? showBorderWidth;
  final double? width;
  final double? height;
  const Avatar({
    this.url,
    this.radius,
    this.showBorder = true,
    this.size = 48,
    this.width,
    this.height,
    super.key,
    this.showBorderColor,
    this.showBorderWidth = 2,
  });

  @override
  Widget build(BuildContext context) {
    if (url == null || (url?.isEmpty ?? true)) {
      return Container(
        width: width ?? size,
        height: height ?? size,
        decoration: BoxDecoration(
          border: (showBorder && showBorderColor != null)
              ? Border.all(
                  color: showBorderColor ?? Colors.grey.withOpacity(.2),
                  width: showBorderWidth ?? 2,
                )
              : null,
          color: Colors.grey.withOpacity(.2),
          borderRadius: radius != null
              ? BorderRadius.circular(radius ?? 16)
              : AppStyle.radius32,
        ),
        child: const Icon(
          Icons.person,
          color: AppColors.greyF6F6F6,
          size: 24,
        ),
      );
    }
    return Container(
      width: width ?? size,
      height: height ?? size,
      decoration: BoxDecoration(
        borderRadius: radius != null
            ? BorderRadius.circular(radius ?? 16)
            : AppStyle.radius32,
        border: (showBorder && showBorderColor != null)
            ? Border.all(
                color: showBorderColor ?? Colors.grey.withOpacity(.2),
                width: showBorderWidth ?? 2,
              )
            : null,
      ),
      child: NetWorkCacheImage(
        url!,
        width: width ?? size,
        height: height ?? size,
        borderRadius: radius != null ? radius! : 16,
      ),
    );
  }
}
