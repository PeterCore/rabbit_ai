import 'package:flutter/material.dart';

import '../app/app_color.dart';

class RoundButtonIcon extends StatelessWidget {
  final String? text;
  final VoidCallback onPressed;
  final Color backgroundColor;
  final Color textColor;
  final double borderRadius;
  final double horizontalPadding;
  final double verticalPadding;
  final TextStyle? textStyle;
  final Size? minSize;
  final double? distance;
  final Widget? icon;
  final BorderSide? side;
  const RoundButtonIcon(
      {super.key,
      this.text,
      this.minSize,
      required this.onPressed,
      this.backgroundColor = AppColors.pinkColor, // 默认粉色
      this.textColor = Colors.white,
      this.borderRadius = 26,
      this.horizontalPadding = 10,
      this.verticalPadding = 4,
      this.textStyle,
      this.distance,
      this.side,
      this.icon});
  @override
  Widget build(BuildContext context) {
    return ElevatedButton(
      style: ElevatedButton.styleFrom(
        shadowColor: Colors.transparent,
        minimumSize: minSize ?? const Size(64, 26),
        // maximumSize: minSize ?? const Size(64, 26),
        foregroundColor: AppColors.whiteColor,
        backgroundColor: backgroundColor,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(borderRadius),
          side: side ?? BorderSide.none,
        ),
        padding: EdgeInsets.symmetric(
          horizontal: horizontalPadding,
          vertical: verticalPadding,
        ),
      ),
      onPressed: () {
        onPressed.call();
      },
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        mainAxisSize: MainAxisSize.min,
        children: [
          icon ?? Container(),
          SizedBox(width: distance ?? 4), // 根据需要调整 icon 与文字的间距
          Text(text ?? "", style: textStyle),
        ],
      ),
    );
  }
}

class RoundedButton extends StatelessWidget {
  final String text;
  final VoidCallback onPressed;
  final Color backgroundColor;
  final Color textColor;
  final double borderRadius;
  final double horizontalPadding;
  final double verticalPadding;
  final TextStyle? textStyle;
  final Size? minSize;
  final BorderSide? side;
  const RoundedButton({
    super.key,
    required this.text,
    this.minSize,
    required this.onPressed,
    this.backgroundColor = AppColors.pinkColor, // 默认粉色
    this.textColor = Colors.white,
    this.borderRadius = 26,
    this.horizontalPadding = 10,
    this.verticalPadding = 4,
    this.textStyle,
    this.side,
  });

  @override
  Widget build(BuildContext context) {
    return ElevatedButton(
      style: ElevatedButton.styleFrom(
        shadowColor: Colors.transparent,
        minimumSize: minSize,
        foregroundColor: textColor,
        backgroundColor: backgroundColor,
        shape: RoundedRectangleBorder(
          side: side ?? BorderSide.none,
          borderRadius: BorderRadius.circular(borderRadius),
        ),
        padding: EdgeInsets.symmetric(
          horizontal: horizontalPadding,
          vertical: verticalPadding,
        ),
      ),
      onPressed: onPressed,
      child: Text(
        text,
        style: textStyle ?? TextStyle(color: textColor),
      ),
    );
  }
}
