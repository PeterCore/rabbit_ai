import 'dart:convert';

import 'package:crypto/crypto.dart';

extension ExtensionString on String {
  bool get isNullOrEmpty => this == null || this!.isEmpty;

  bool get isNotNullOrEmpty => !isNullOrEmpty;

  bool get isValidCNMobile {
    final regex = RegExp(r'^1[3-9]\d{9}$');
    return regex.hasMatch(this);
  }

  int get age {
    final birth = DateTime.parse(this);
    final now = DateTime.now();
    var years = now.year - birth.year;
    if (now.month < birth.month ||
        (now.month == birth.month && now.day < birth.day)) {
      years--;
    }
    return years;
  }

  String truncateText(int maxLength) {
    if (length <= maxLength) return this;
    return '${substring(0, maxLength)}...';
  }

  String maskPhone() {
    if (length < 10) return this;
    return replaceRange(6, length - 4, '****');
  }

  String maskXxPhone() {
    if (length != 11) return this;
    return replaceRange(3, length - 4, '*xx*');
  }

  String toMd5() {
    var bytes = utf8.encode(this);
    var digest = md5.convert(bytes);
    return digest.toString();
  }

  bool isOnlyZhEn() {
    final reg = RegExp(r'^[A-Za-z\u4E00-\u9FFF]+$');
    return reg.hasMatch(this);
  }

  int stringRunsLength(String s) {
    int len = 0;
    for (var r in s.runes) {
      if (r > 0xFF) {
        len += 2;
      } else {
        len += 1;
      }
    }
    return len;
  }

  bool isLengthOK({int maxLen = 12}) => stringRunsLength(this) <= maxLen;
}
