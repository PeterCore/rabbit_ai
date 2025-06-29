import 'package:flutter/material.dart';

import 'app_color.dart';

class AppTextStyles {
  static textWordStyle(
      {Color? txtColor = AppColors.blackColor,
      double fontSize = 14,
      FontWeight fontWeight = FontWeight.w400,
      TextDecoration decoration = TextDecoration.none,
      double? letterSpacing,
      double? height,
      TextOverflow? overflow}) {
    return TextStyle(
        fontFamily: 'PingFang SC',
        fontWeight: fontWeight,
        color: txtColor,
        overflow: overflow,
        decoration: decoration,
        height: height,
        decorationColor: AppColors.blackColor,
        letterSpacing: letterSpacing,
        fontSize: fontSize);
  }

  static textNumStyle(
      {Color? txtColor = AppColors.blackColor,
      double fontSize = 14,
      FontWeight fontWeight = FontWeight.w900,
      TextDecoration decoration = TextDecoration.none,
      double? letterSpacing,
      double? height,
      TextOverflow? overflow}) {
    return TextStyle(
        fontFamily: 'AlimamaShuHeiTi',
        fontWeight: fontWeight,
        color: txtColor,
        overflow: overflow,
        decoration: decoration,
        height: height,
        decorationColor: AppColors.blackColor,
        letterSpacing: letterSpacing ?? 2,
        fontSize: fontSize);
  }

  static textRobotoStyle(
      {Color? txtColor = AppColors.blackColor,
      double fontSize = 14,
      FontWeight fontWeight = FontWeight.w900,
      TextDecoration decoration = TextDecoration.none,
      double? letterSpacing,
      double? height,
      TextOverflow? overflow}) {
    return TextStyle(
        fontFamily: 'Roboto',
        fontWeight: fontWeight,
        color: txtColor,
        overflow: overflow,
        decoration: decoration,
        height: height,
        decorationColor: AppColors.blackColor,
        letterSpacing: letterSpacing ?? 2,
        fontSize: fontSize);
  }
}
