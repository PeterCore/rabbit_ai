import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';
import 'package:rabbit_ai_fe/app/app_dialog_utils.dart';
import 'package:rabbit_ai_fe/common/constants/size_extensions.dart';
import 'package:rabbit_ai_fe/common/screen_util.dart';
import '../app/app_color.dart';
import '../app/app_text_styles.dart';
import '../common/constants/size_constants.dart';

Future<void> copyToClipboard(String text) async {
  await Clipboard.setData(ClipboardData(text: text));
}

class PositiveIntegerFormatter extends TextInputFormatter {
  static final _regExp = RegExp(r'^[1-9]\d*$');

  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    final text = newValue.text;
    if (text.isEmpty) {
      return newValue;
    }
    if (_regExp.hasMatch(text)) {
      return newValue;
    }
    return oldValue;
  }
}

void showRichTextDialog(BuildContext context,
    {String prefix = "",
    String mid = "",
    String suffix = "",
    Function(bool? value)? onClick}) {
  AppDialogUtils.showContentAlertDialog(
    context,
    RichText(
      textAlign: TextAlign.center,
      text: TextSpan(children: [
        TextSpan(
            text: prefix,
            style: AppTextStyles.textWordStyle(
                txtColor: AppColors.blackColorO6,
                fontWeight: FontWeight.w500,
                fontSize: Sizes.dimen_12.sp)),
        TextSpan(
            text: mid,
            style: AppTextStyles.textWordStyle(
                fontWeight: FontWeight.w500,
                txtColor: AppColors.pinkColor,
                fontSize: Sizes.dimen_12.sp)),
        TextSpan(
            text: suffix,
            style: AppTextStyles.textWordStyle(
                fontWeight: FontWeight.w500,
                txtColor: AppColors.blackColorO6,
                fontSize: Sizes.dimen_12.sp)),
      ]),
    ),
    title: "提示",
  ).then((value) {
    onClick?.call(value);
    // if (value == true) {
    //   DialogUtils.showToast("已发送申请");
    // }
  });
}

Widget commonListViewSaparated({
  EdgeInsetsGeometry? padding,
  ScrollPhysics? physics,
  ScrollController? controller,
  int itemCount = 0,
  Widget Function(BuildContext, int)? separatorBuilder,
  required Widget? Function(BuildContext, int) itemBuilder,
}) {
  return ListView.separated(
    shrinkWrap: true,
    padding: padding,
    physics: physics ?? const ClampingScrollPhysics(),
    controller: controller,
    itemCount: itemCount,
    separatorBuilder: separatorBuilder ?? (ctx, index) => const SizedBox(),
    itemBuilder: itemBuilder,
  );
}

Widget underlineTextField(
    {Color textColor = AppColors.blackColor,
    Color hintTextColor = AppColors.blackColor,
    double? fontSize,
    FontWeight? fontWeight,
    TextEditingController? controller,
    OutlineInputBorder? border,
    Color backgroundColor = Colors.white,
    Color counterColor = Colors.white,
    Color? underLineColor,
    Widget? prefix,
    bool? obscureText,
    String? initialValue,
    Widget? suffixWidget,
    bool suffix = false,
    FormFieldValidator? validator,
    bool autoFocus = false,
    bool readOnly = false,
    bool enableInteractiveSelection = true,
    List<TextInputFormatter>? inputFormatter,
    TextInputType? keyBoardType,
    String hintText = '',
    double borderRadius = 8,
    TextInputAction? textInputAction,
    bool? enabled,
    Color? focusColor,
    Color? cursorColor,
    double prefixMaxWidth = 34,
    FocusNode? focusNode,
    int? maxLines,
    InputBorder? disabledBorder,
    int? maxLength,
    FontWeight? hintStyleFontWeight,
    Function()? onTap,
    Function(String val)? onChanged,
    Function(String val)? onSubmitted,
    double verticalPadding = 0,
    double? horizontalPadding,
    List<String>? autofillHints}) {
  return TextFormField(
    initialValue: initialValue,
    enabled: enabled,
    autofillHints: autofillHints,
    cursorColor: cursorColor ?? AppColors.pinkColor,
    autovalidateMode: AutovalidateMode.onUserInteraction,
    controller: controller,
    focusNode: focusNode,
    enableInteractiveSelection: enableInteractiveSelection,
    validator: validator,
    readOnly: readOnly,
    textInputAction: textInputAction,
    onTap: onTap,
    onChanged: onChanged,
    autofocus: autoFocus,
    textCapitalization: TextCapitalization.none,
    keyboardType: keyBoardType ?? TextInputType.visiblePassword,
    obscureText: obscureText ?? false,
    maxLines: maxLines ?? 1,
    inputFormatters: inputFormatter ?? [],
    maxLength: maxLength,
    maxLengthEnforcement: MaxLengthEnforcement.none,
    onFieldSubmitted: onSubmitted,
    decoration: InputDecoration(
      hintText: hintText,
      counterStyle: AppTextStyles.textWordStyle(txtColor: counterColor),
      filled: true,
      fillColor: backgroundColor,
      errorMaxLines: 5,
      disabledBorder: UnderlineInputBorder(
        borderSide: BorderSide(
          color: underLineColor ?? AppColors.greyF8F8F8,
          width: 1,
        ),
      ),
      enabledBorder: UnderlineInputBorder(
        borderSide: BorderSide(
          color: underLineColor ?? AppColors.greyF8F8F8,
          width: 1,
        ),
      ),
      focusedBorder: UnderlineInputBorder(
        borderSide: BorderSide(
          color: underLineColor ?? AppColors.greyF8F8F8,
          width: 1,
        ),
      ),
      border: UnderlineInputBorder(
        borderSide: BorderSide(
          color: underLineColor ?? AppColors.greyF8F8F8,
          width: 1,
        ),
      ),
      alignLabelWithHint: false,
      hintStyle: AppTextStyles.textWordStyle(
          fontWeight: hintStyleFontWeight ?? FontWeight.w400,
          txtColor: hintTextColor,
          fontSize: fontSize ?? 14),
      prefixIcon: prefix != null
          ? Padding(
              padding: EdgeInsets.only(right: horizontalPadding ?? 10),
              child: prefix,
            )
          : null,
      contentPadding: prefix == null
          ? EdgeInsets.symmetric(
              horizontal: horizontalPadding ?? 10, vertical: verticalPadding)
          : null,
      prefixIconConstraints: BoxConstraints(maxWidth: prefixMaxWidth),
    ),
    style: AppTextStyles.textWordStyle(
        txtColor: textColor,
        fontSize: fontSize ?? 14,
        fontWeight: fontWeight ?? FontWeight.w400),
  );
}

// String fromName = "",

Widget commonTextField(
    {Color textColor = AppColors.blackColor,
    Color hintTextColor = AppColors.blackColor,
    double? fontSize,
    TextEditingController? controller,
    OutlineInputBorder? border,
    Color backgroundColor = Colors.white,
    Color counterColor = Colors.white,
    Widget? prefix,
    bool? obscureText,
    String? initialValue,
    Widget? suffixWidget,
    FontWeight? fontWeight,
    TextInputAction? textInputAction,
    bool suffix = false,
    FormFieldValidator? validator,
    bool autoFocus = false,
    bool readOnly = false,
    bool enableInteractiveSelection = true,
    List<TextInputFormatter>? inputFormatter,
    TextInputType? keyBoardType,
    String hintText = '',
    double borderRadius = 8,
    bool? enabled,
    Color? focusColor,
    Color? cursorColor,
    double prefixMaxWidth = 34,
    FocusNode? focusNode,
    int? maxLines,
    InputBorder? disabledBorder,
    int? maxLength,
    Function()? onTap,
    Function(String val)? onChanged,
    Function(String val)? onSubmitted,
    double verticalPadding = 0,
    double? horizontalPadding,
    List<String>? autofillHints}) {
  return TextFormField(
    initialValue: initialValue,
    textInputAction: textInputAction,
    enabled: enabled,
    autofillHints: autofillHints,
    cursorColor: cursorColor ?? AppColors.pinkColor,
    autovalidateMode: AutovalidateMode.onUserInteraction,
    controller: controller,
    focusNode: focusNode,
    enableInteractiveSelection: enableInteractiveSelection,
    validator: validator,
    readOnly: readOnly,
    onTap: onTap,
    onChanged: onChanged,
    autofocus: autoFocus,
    textCapitalization: TextCapitalization.none,
    keyboardType: keyBoardType ?? TextInputType.visiblePassword,
    obscureText: obscureText ?? false,
    maxLines: maxLines ?? 1,
    inputFormatters: inputFormatter ?? [],
    maxLength: maxLength,
    onFieldSubmitted: onSubmitted,
    maxLengthEnforcement: MaxLengthEnforcement.none,
    decoration: InputDecoration(
      hintText: hintText,
      counterStyle: AppTextStyles.textWordStyle(txtColor: counterColor),
      filled: true,
      fillColor: backgroundColor,
      errorMaxLines: 5,
      disabledBorder: disabledBorder ??
          OutlineInputBorder(
              borderSide: const BorderSide(
                color: Colors.transparent,
                width: 1,
              ),
              borderRadius: BorderRadius.circular(borderRadius)),
      enabledBorder: border ??
          OutlineInputBorder(
              borderSide: const BorderSide(
                color: Colors.transparent,
                width: 1,
              ),
              borderRadius: BorderRadius.circular(borderRadius)),
      focusedBorder: border ??
          OutlineInputBorder(
              borderSide: BorderSide(
                color: focusColor ?? Colors.transparent,
                width: 1,
              ),
              borderRadius: BorderRadius.circular(borderRadius)),
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(borderRadius),
      ),
      alignLabelWithHint: false,
      hintStyle: AppTextStyles.textWordStyle(
          fontWeight: FontWeight.w400,
          txtColor: hintTextColor,
          fontSize: fontSize ?? 14),
      prefixIcon: prefix,
      contentPadding: EdgeInsets.symmetric(
          horizontal: horizontalPadding ?? 10, vertical: verticalPadding),
      prefixIconConstraints: BoxConstraints(maxWidth: prefixMaxWidth),
    ),
    style: AppTextStyles.textWordStyle(
        txtColor: textColor,
        fontSize: fontSize ?? 14,
        fontWeight: fontWeight ?? FontWeight.w400),
  );
}

Widget materialButtonWithChild(
    {Widget? child,
    Color color = AppColors.blackColor,
    BorderRadius? borderRadiusOnly,
    double borderRadius = 8,
    Color borderColor = Colors.transparent,
    double? borderWidth,
    double height = 0,
    double width = 0,
    double elevation = 0,
    Color splashColor = AppColors.grey9493AC,
    Color? highlightColor,
    // double highlightElevation = 3,
    bool shadow = false,
    EdgeInsets padding = const EdgeInsets.symmetric(horizontal: 6, vertical: 3),
    void Function()? onPressed}) {
  return Material(
    color: Colors.white,
    child: Container(
        decoration: BoxDecoration(
          color: AppColors.whiteColor,
          borderRadius: BorderRadius.circular(borderRadius),
          boxShadow: shadow
              ? <BoxShadow>[
                  BoxShadow(
                    color: AppColors.greyE0E0E0.withOpacity(0.1),
                    blurRadius: borderRadius,
                    offset: const Offset(0, 10),
                  ),
                  BoxShadow(
                    color: AppColors.greyE0E0E0,
                    blurRadius: borderRadius,
                    offset: const Offset(10, 0),
                  ),
                  BoxShadow(
                    color: AppColors.greyE0E0E0.withOpacity(0.1),
                    blurRadius: borderRadius,
                    offset: const Offset(0, -10),
                  ),
                  BoxShadow(
                    color: AppColors.greyE0E0E0,
                    blurRadius: borderRadius,
                    offset: const Offset(-10, 0),
                  ),
                ]
              : <BoxShadow>[],
        ),
        child: MaterialButton(
            materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
            padding: padding,
            elevation: elevation,
            highlightElevation: 0,
            splashColor: splashColor,
            highlightColor: highlightColor,
            clipBehavior: Clip.hardEdge,
            shape: OutlineInputBorder(
                borderRadius:
                    borderRadiusOnly ?? BorderRadius.circular(borderRadius),
                borderSide: borderWidth == null
                    ? BorderSide.none
                    : BorderSide(color: borderColor, width: borderWidth)),
            height: height,
            minWidth: width,
            color: color,
            onPressed: onPressed ?? () {},
            child: child)),
  );
}
