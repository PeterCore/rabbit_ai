import 'package:rabbit_ai_fe/widgets/common_codes.dart';
import 'package:rabbit_ai_fe/widgets/round_button.dart';
import 'package:rabbit_ai_fe/widgets/round_container.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';

import 'app_color.dart';
import 'app_style.dart';
import 'app_text_styles.dart';

// Future<LottieComposition> loadComposition() async {
//   final bytes = await rootBundle.load('assets/animation/data.json');
//   return await LottieComposition.fromByteData(bytes);
// }

class CustomEditWidget extends StatefulWidget {
  const CustomEditWidget({super.key, this.customs = const [], this.onSelected});
  final List<String> customs;
  final Function(String title)? onSelected;
  @override
  State<CustomEditWidget> createState() => _CustomEditWidgetState();
}

class _CustomEditWidgetState extends State<CustomEditWidget> {
  final textEditingController = TextEditingController();
  String selectTitle = "";
  @override
  void initState() {
    // TODO: implement initState
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: MediaQuery.of(context).size.width - 100,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          ...widget.customs.map((item) => Padding(
                padding: const EdgeInsets.only(bottom: 10),
                child: SizedBox(
                  width: MediaQuery.of(context).size.width - 100 - 40,
                  child: RoundedButton(
                      verticalPadding: 4,
                      horizontalPadding: 80,
                      backgroundColor: selectTitle == item
                          ? AppColors.pinkBgColor
                          : AppColors.whiteColor,
                      textStyle: AppTextStyles.textWordStyle(
                          txtColor: selectTitle == item
                              ? AppColors.pinkColor
                              : AppColors.greyAEAEAE,
                          fontWeight: FontWeight.w500,
                          fontSize: 16),
                      side: selectTitle != item
                          ? const BorderSide(
                              color: AppColors.greyAEAEAE,
                              width: 1,
                            )
                          : BorderSide.none,
                      minSize: const Size(100, 36),
                      text: item,
                      onPressed: () {
                        widget.onSelected?.call(item);
                        setState(() {
                          selectTitle = item;
                        });
                      }),
                ),
              )),
          Row(
            children: [
              Flexible(
                flex: 3,
                child: Text(
                  "自定义板块",
                  style: AppTextStyles.textWordStyle(
                      txtColor: AppColors.blackColor,
                      fontWeight: FontWeight.w500,
                      fontSize: 16),
                ),
              ),
              const SizedBox(
                width: 14,
              ),
              Flexible(
                flex: 7,
                child: SizedBox(
                  height: 44,
                  child: commonTextField(
                    borderRadius: 22,
                    keyBoardType: TextInputType.text,
                    hintText: "请输入你想添加的板块",
                    controller: textEditingController,
                    hintTextColor: AppColors.blackColorO4,
                    fontSize: 14,
                    backgroundColor: AppColors.greyF8F8F8,
                    cursorColor: AppColors.pinkColor,
                    autoFocus: false,
                    onChanged: (val) {},
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(
            height: 14,
          ),
          SizedBox(
            height: 44,
            width: double.infinity,
            child: RoundedButton(
                horizontalPadding: 44,
                verticalPadding: 10,
                text: "确定",
                textStyle: AppTextStyles.textWordStyle(
                    txtColor: AppColors.whiteColor,
                    fontWeight: FontWeight.w500,
                    fontSize: 14),
                backgroundColor: AppColors.pinkColor,
                onPressed: (() {
                  widget.onSelected?.call(textEditingController.text);
                  // Navigator.of(context).pop("");
                  //onCallBack?.call(textEditingController.text);
                  // Navigator.of(context).pop(textEditingController.text);
                })),
          ),
        ],
      ),
    );
  }
}

class AppDialogUtils {
  // static void showRibbons() {
  //   SmartDialog.showToast('',
  //       displayType: SmartToastType.multi,
  //       alignment: Alignment.center,
  //       maskColor: AppColors.greyAEAEAE,
  //       builder: (_) => const RibbonsView());
  // }

  static void showToast(String title, {Alignment? alignment}) {
    SmartDialog.showToast(
      '',
      displayType: SmartToastType.multi,
      alignment: alignment ?? Alignment.center,
      maskColor: AppColors.greyAEAEAE,
      builder: (_) => RoundContainer(
        backgroundColor: Colors.grey.withOpacity(0.8),
        boxShadow: const [
          BoxShadow(
              color: Colors.transparent,
              offset: Offset(0.0, 0.0), //阴影y轴偏移量
              blurRadius: 0, //阴影模糊程度
              spreadRadius: 0.5 //阴影扩散程度
              ),
        ],
        horizontal: 12,
        vertical: 12,
        child: Text(
          title,
          style: AppTextStyles.textWordStyle(txtColor: Colors.white),
        ),
      ),
    );
    // SmartDialog.showToast(title,
    //     alignment: alignment ?? Alignment.center,
    //     maskColor: AppColors.greyAEAEAE);
  }

  static Future<bool?> showContentAlertDialog(
    BuildContext context,
    Widget content, {
    String title = '',
    String confirm = '',
    String cancel = '',
    bool barrierDismissible = true,
    List<Widget>? actions,
  }) async {
    final alterDialog = AlertDialog(
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(24.0),
      ),
      titlePadding: const EdgeInsets.only(top: 18),
      contentPadding: EdgeInsets.zero,
      actionsPadding: const EdgeInsets.only(bottom: 18),
      title: Text(
        title,
        textAlign: TextAlign.center,
        style: AppTextStyles.textWordStyle(
            txtColor: AppColors.blackColor,
            fontWeight: FontWeight.w500,
            fontSize: 18),
      ),
      content: Padding(padding: AppStyle.edgeInsetsV12, child: content),
      actionsAlignment: MainAxisAlignment.center,
      actions: [
        RoundedButton(
            horizontalPadding: 40,
            verticalPadding: 10,
            text: cancel.isEmpty ? "取消" : cancel,
            textStyle: AppTextStyles.textWordStyle(
                txtColor: AppColors.pinkColor,
                fontWeight: FontWeight.w500,
                fontSize: 14),
            backgroundColor: AppColors.pinkBgColor,
            onPressed: (() => Navigator.of(context).pop(false))),
        RoundedButton(
            horizontalPadding: 40,
            verticalPadding: 10,
            text: confirm.isEmpty ? "确定" : confirm,
            textStyle: AppTextStyles.textWordStyle(
                txtColor: AppColors.whiteColor,
                fontWeight: FontWeight.w500,
                fontSize: 14),
            backgroundColor: AppColors.pinkColor,
            onPressed: (() => Navigator.of(context).pop(true))),
        ...?actions,
      ],
    );

    return showDialog(
      context: context,
      barrierDismissible: barrierDismissible,
      builder: (BuildContext context) {
        return alterDialog;
      },
    );
  }

  static Future<bool?> showAlertDialog(
    BuildContext context,
    String content, {
    String title = '',
    String confirm = '',
    String cancel = '',
    bool barrierDismissible = true,
    List<Widget>? actions,
  }) async {
    final alterDialog = AlertDialog(
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(24.0),
      ),
      titlePadding: const EdgeInsets.only(top: 18),
      contentPadding: EdgeInsets.zero,
      actionsPadding: const EdgeInsets.only(bottom: 18),
      title: Text(
        title,
        textAlign: TextAlign.center,
        style: AppTextStyles.textWordStyle(
            txtColor: AppColors.blackColor,
            fontWeight: FontWeight.w500,
            fontSize: 18),
      ),
      content: content.isNotEmpty
          ? Padding(
              padding: AppStyle.edgeInsetsV12,
              child: Text(
                content,
                textAlign: TextAlign.center,
                style: AppTextStyles.textWordStyle(
                    txtColor: AppColors.blackColorO6,
                    fontWeight: FontWeight.w400,
                    fontSize: 12),
              ),
            )
          : null,
      actionsAlignment: MainAxisAlignment.center,
      actions: [
        RoundedButton(
            horizontalPadding: 40,
            verticalPadding: 10,
            text: cancel.isEmpty ? "取消" : cancel,
            textStyle: AppTextStyles.textWordStyle(
                txtColor: AppColors.pinkColor,
                fontWeight: FontWeight.w500,
                fontSize: 14),
            backgroundColor: AppColors.pinkBgColor,
            onPressed: (() => Navigator.of(context).pop(false))),
        RoundedButton(
            horizontalPadding: 40,
            verticalPadding: 10,
            text: confirm.isEmpty ? "确定" : confirm,
            textStyle: AppTextStyles.textWordStyle(
                txtColor: AppColors.whiteColor,
                fontWeight: FontWeight.w500,
                fontSize: 14),
            backgroundColor: AppColors.pinkColor,
            onPressed: (() => Navigator.of(context).pop(true))),
        ...?actions,
      ],
    );

    return showDialog(
      context: context,
      barrierColor: Colors.grey.withOpacity(0.5), // 灰色半透明遮罩
      barrierDismissible: barrierDismissible,
      builder: (BuildContext context) {
        return alterDialog;
      },
    );
  }

  static Future<bool?> showMessageDialog(BuildContext context, String content,
      {String title = '', String confirm = '', bool selectable = false}) async {
    final alterDialog = AlertDialog(
      title: Text(title),
      content: Padding(
        padding: AppStyle.edgeInsetsV12,
        child: selectable ? SelectableText(content) : Text(content),
      ),
      actions: [
        TextButton(
          onPressed: (() => Navigator.of(context).pop(true)),
          child: Text(confirm.isEmpty ? "确定" : confirm),
        ),
      ],
    );
    return showDialog(
      context: context,
      builder: (BuildContext context) {
        return alterDialog;
      },
    );
  }

  static Future<String?> showCustomEditDialog(BuildContext context,
      {String title = "",
      String confirm = "",
      List<String> customs = const []}) async {
    final alterDialog = AlertDialog(
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(24.0),
      ),
      insetPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 12),
      actionsPadding: const EdgeInsets.only(bottom: 18, left: 30, right: 30),
      titlePadding: const EdgeInsets.only(top: 12),
      title: Text(
        title,
        textAlign: TextAlign.center,
        style: AppTextStyles.textWordStyle(
            txtColor: AppColors.blackColor,
            fontWeight: FontWeight.w500,
            fontSize: 24),
      ),
      content: CustomEditWidget(
        customs: customs,
        onSelected: (value) {
          Navigator.of(context).pop(value);
        },
      ),
      actionsAlignment: MainAxisAlignment.center,
      // actions: [
      //   SizedBox(
      //     height: 44,
      //     width: double.infinity,
      //     child: RoundedButton(
      //         horizontalPadding: 44,
      //         verticalPadding: 10,
      //         text: confirm.isEmpty ? "确定" : confirm,
      //         textStyle: AppTextStyles.textWordStyle(
      //             txtColor: AppColors.whiteColor,
      //             fontWeight: FontWeight.w500,
      //             fontSize: 14),
      //         backgroundColor: AppColors.pinkColor,
      //         onPressed: (() {
      //           Navigator.of(context).pop("");
      //           //onCallBack?.call(textEditingController.text);
      //           // Navigator.of(context).pop(textEditingController.text);
      //         })),
      //   ),
      // ],
    );

    return showDialog(
      context: context,
      barrierColor: Colors.grey.withOpacity(0.5), // 灰色半透明遮罩
      builder: (BuildContext context) {
        return MediaQuery(
            // 使用当前 MediaQuery，但将 viewInsets 固定为 EdgeInsets.zero
            data: MediaQuery.of(context)
                .copyWith(viewInsets: const EdgeInsets.only(bottom: 109)),
            child: alterDialog);
      },
    );
  }

  static Future<String?> showEditTextDialog(
      BuildContext context, String content,
      {String title = '',
      String? hintText,
      String? text,
      String confirm = '',
      Widget? extWidget,
      List<TextInputFormatter>? inputFormatter,
      TextInputType? keyBoardType,
      String cancel = ''}) async {
    final TextEditingController textEditingController =
        TextEditingController(text: text ?? "");
    // final alterDialog = ;

    return showDialog(
      context: context,
      barrierColor: Colors.grey.withOpacity(0.5), // 灰色半透明遮罩
      builder: (BuildContext dialogContext) {
        return MediaQuery(
            // 使用当前 MediaQuery，但将 viewInsets 固定为 EdgeInsets.zero
            data: MediaQuery.of(context)
                .copyWith(viewInsets: const EdgeInsets.only(bottom: 36)),
            child: AlertDialog(
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(24.0),
              ),
              actionsPadding:
                  const EdgeInsets.only(bottom: 18, left: 12, right: 12),
              title: Text(
                title,
                textAlign: TextAlign.center,
                style: AppTextStyles.textWordStyle(
                    txtColor: AppColors.blackColor,
                    fontWeight: FontWeight.w500,
                    fontSize: 18),
              ),
              content: SizedBox(
                height: 36,
                child: Row(children: [
                  Expanded(
                    child: commonTextField(
                      borderRadius: 18,
                      keyBoardType: keyBoardType ?? TextInputType.text,
                      hintText: "请输入",
                      controller: textEditingController,
                      hintTextColor: AppColors.blackColorO4,
                      fontSize: 14,
                      inputFormatter: inputFormatter,
                      backgroundColor: AppColors.greyF8F8F8,
                      cursorColor: AppColors.pinkColor,
                      autoFocus: false,
                      onChanged: (val) {},
                    ),
                  ),
                  extWidget ?? Container()
                ]),
              ),
              actionsAlignment: MainAxisAlignment.center,
              actions: [
                SizedBox(
                  height: 40,
                  width: double.infinity,
                  child: RoundedButton(
                      horizontalPadding: 40,
                      verticalPadding: 10,
                      text: confirm.isEmpty ? "确定" : cancel,
                      textStyle: AppTextStyles.textWordStyle(
                          txtColor: AppColors.whiteColor,
                          fontWeight: FontWeight.w500,
                          fontSize: 14),
                      backgroundColor: AppColors.pinkColor,
                      onPressed: (() {
                        var text = textEditingController.text;
                        Navigator.of(dialogContext).pop(text);
                      })),
                ),
              ],
            ));
      },
    );
  }

  static showContentEditTextDialog(BuildContext context,
      {String title = '',
      String? hintText,
      Widget? content,
      String confirm = '',
      Function(String content)? onCallBack,
      String cancel = ''}) async {
    final TextEditingController textEditingController =
        TextEditingController(text: "");
    final alterDialog = AlertDialog(
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(24.0),
      ),
      actionsPadding: const EdgeInsets.only(bottom: 18, left: 12, right: 12),
      title: Text(
        title,
        textAlign: TextAlign.center,
        style: AppTextStyles.textWordStyle(
            txtColor: AppColors.blackColor,
            fontWeight: FontWeight.w500,
            fontSize: 18),
      ),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          content ?? Container(),
          AppStyle.vGap12,
          SizedBox(
            height: 44,
            child: commonTextField(
              borderRadius: 10,
              keyBoardType: TextInputType.text,
              hintText: "请输入",
              controller: textEditingController,
              hintTextColor: AppColors.blackColorO4,
              fontSize: 14,
              backgroundColor: AppColors.greyF8F8F8,
              cursorColor: AppColors.pinkColor,
              autoFocus: false,
              onChanged: (val) {},
            ),
          ),
        ],
      ),
      actionsAlignment: MainAxisAlignment.center,
      actions: [
        SizedBox(
          height: 40,
          width: double.infinity,
          child: RoundedButton(
              horizontalPadding: 40,
              verticalPadding: 10,
              text: confirm.isEmpty ? "确定" : cancel,
              textStyle: AppTextStyles.textWordStyle(
                  txtColor: AppColors.whiteColor,
                  fontWeight: FontWeight.w500,
                  fontSize: 14),
              backgroundColor: AppColors.pinkColor,
              onPressed: (() {
                onCallBack?.call(textEditingController.text);
                // Navigator.of(context).pop(textEditingController.text);
              })),
        ),
      ],
    );

    return showDialog(
      context: context,
      barrierColor: Colors.grey.withOpacity(0.5), // 灰色半透明遮罩
      builder: (BuildContext context) {
        return MediaQuery(
            // 使用当前 MediaQuery，但将 viewInsets 固定为 EdgeInsets.zero
            data: MediaQuery.of(context)
                .copyWith(viewInsets: const EdgeInsets.only(bottom: 36)),
            child: alterDialog);
      },
    );
  }
}
