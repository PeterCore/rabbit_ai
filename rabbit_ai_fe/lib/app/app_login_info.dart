import 'dart:async';

import 'package:flutter/cupertino.dart';

import '../common/storage_uitl.dart';

final tokenNotifier = ValueNotifier<String?>(null);

class AppLoginInfo {
  factory AppLoginInfo() => _getInstance()!;

  static AppLoginInfo? get instance => _getInstance();
  static AppLoginInfo? _instance;

  AppLoginInfo._internal();

  String _appToken = "";
  String get appToken => _appToken;
  bool _isTeaParty = false;
  bool get isTeaParty => _isTeaParty;
  static AppLoginInfo? _getInstance() {
    _instance ??= AppLoginInfo._internal();
    return _instance;
  }

  Future saveTeaParty(bool isTeaParty) async {
    _isTeaParty = isTeaParty;
    await StorageUtils.setBoolValue("isTeaParty", _isTeaParty);
  }

  Future saveAppToken(String token) async {
    _appToken = token;
    await StorageUtils.setValue("app_token", token);
  }

  Future initAppToken() async {
    _appToken = await StorageUtils.getValue("app_token");
    _isTeaParty = await StorageUtils.getBoolValue("isTeaParty");
  }

  Future clearToken() async {
    _appToken = "";
    await StorageUtils.setValue("app_token", "");
  }
}
