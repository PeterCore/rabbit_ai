import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'app_color.dart';

class AppStyle {
  static ThemeData lightTheme = ThemeData.light(
    useMaterial3: false,
  ).copyWith(
    brightness: Brightness.light,
    colorScheme: AppColors.colorSchemeLight,
    scaffoldBackgroundColor: Colors.white,
    cardColor: Colors.white,
    appBarTheme: AppBarTheme(
      elevation: 0,
      backgroundColor: Colors.transparent,
      foregroundColor: Colors.white,
      centerTitle: false,
      shape: Border(
        bottom: BorderSide(
          color: Colors.grey.withOpacity(.2),
          width: 1,
        ),
      ),
      iconTheme: const IconThemeData(
        color: AppColors.black333,
      ),
      titleTextStyle: const TextStyle(
        fontSize: 16,
        color: AppColors.black333,
      ),
      systemOverlayStyle: SystemUiOverlayStyle.dark.copyWith(
        systemNavigationBarColor: AppColors.whiteColor, // 设置底部导航栏颜色
      ),
    ),
  );
  static ThemeData darkTheme = ThemeData.dark(
    useMaterial3: false,
  ).copyWith(
    brightness: Brightness.dark,
    primaryColor: Colors.blue,
    cardColor: const Color(0xff424242),
    colorScheme: AppColors.colorSchemeDark,
    scaffoldBackgroundColor: Colors.black,
    tabBarTheme: const TabBarTheme(
      indicatorColor: Colors.blue,
    ),
    appBarTheme: AppBarTheme(
      elevation: 0,
      backgroundColor: Colors.transparent,
      foregroundColor: Colors.white,
      centerTitle: false,
      shape: Border(
        bottom: BorderSide(
          color: Colors.grey.withOpacity(.2),
          width: 1,
        ),
      ),
      titleTextStyle: const TextStyle(
        fontSize: 16,
        color: Colors.white,
      ),
      iconTheme: const IconThemeData(
        color: Colors.white,
      ),
      systemOverlayStyle: SystemUiOverlayStyle.light.copyWith(
        systemNavigationBarColor: Colors.transparent,
      ),
    ),
  );
  static const vGap2 = SizedBox(
    height: 2,
  );
  static const vGap4 = SizedBox(
    height: 4,
  );
  static const vGap6 = SizedBox(
    height: 6,
  );
  static const vGap8 = SizedBox(
    height: 8,
  );
  static const vGap10 = SizedBox(
    height: 10,
  );
  static const vGap11 = SizedBox(
    height: 11,
  );
  static const vGap12 = SizedBox(
    height: 12,
  );
  static const vGap16 = SizedBox(
    height: 16,
  );
  static const vGap18 = SizedBox(
    height: 18,
  );
  static const vGap20 = SizedBox(
    height: 20,
  );
  static const vGap24 = SizedBox(
    height: 24,
  );
  static const vGap29 = SizedBox(
    height: 29,
  );
  static const vGap32 = SizedBox(
    height: 32,
  );
  static const vGap36 = SizedBox(
    height: 36,
  );
  static const vGap41 = SizedBox(
    height: 41,
  );
  static const vGap44 = SizedBox(
    height: 44,
  );
  static const vGap56 = SizedBox(
    height: 56,
  );
  static const vGap70 = SizedBox(
    height: 70,
  );
  static const hGap2 = SizedBox(
    width: 2,
  );
  static const hGap4 = SizedBox(
    width: 4,
  );
  static const hGap6 = SizedBox(
    width: 6,
  );
  static const hGap8 = SizedBox(
    width: 8,
  );
  static const hGap10 = SizedBox(
    width: 10,
  );
  static const hGap12 = SizedBox(
    width: 12,
  );
  static const hGap16 = SizedBox(
    width: 16,
  );
  static const hGap18 = SizedBox(
    width: 18,
  );
  static const hGap24 = SizedBox(
    width: 24,
  );
  static const hGap32 = SizedBox(
    width: 32,
  );

  static const edgeInsetsH4 = EdgeInsets.symmetric(horizontal: 4);
  static const edgeInsetsH8 = EdgeInsets.symmetric(horizontal: 8);
  static const edgeInsetsH12 = EdgeInsets.symmetric(horizontal: 12);
  static const edgeInsetsH16 = EdgeInsets.symmetric(horizontal: 16);
  static const edgeInsetsH18 = EdgeInsets.symmetric(horizontal: 18);
  static const edgeInsetsH20 = EdgeInsets.symmetric(horizontal: 20);
  static const edgeInsetsH24 = EdgeInsets.symmetric(horizontal: 24);

  static const edgeInsetsV4 = EdgeInsets.symmetric(vertical: 4);
  static const edgeInsetsV8 = EdgeInsets.symmetric(vertical: 8);
  static const edgeInsetsV12 = EdgeInsets.symmetric(vertical: 12);
  static const edgeInsetsV24 = EdgeInsets.symmetric(vertical: 24);

  static const edgeInsetsA4 = EdgeInsets.all(4);
  static const edgeInsetsA8 = EdgeInsets.all(8);
  static const edgeInsetsA12 = EdgeInsets.all(12);
  static const edgeInsetsA24 = EdgeInsets.all(24);

  static const edgeInsetsR4 = EdgeInsets.only(right: 4);
  static const edgeInsetsR8 = EdgeInsets.only(right: 8);
  static const edgeInsetsR12 = EdgeInsets.only(right: 12);
  static const edgeInsetsR20 = EdgeInsets.only(right: 20);
  static const edgeInsetsR24 = EdgeInsets.only(right: 24);

  static const edgeInsetsL4 = EdgeInsets.only(left: 4);
  static const edgeInsetsL8 = EdgeInsets.only(left: 8);
  static const edgeInsetsL12 = EdgeInsets.only(left: 12);
  static const edgeInsetsL24 = EdgeInsets.only(left: 24);

  static const edgeInsetsT4 = EdgeInsets.only(top: 4);
  static const edgeInsetsT8 = EdgeInsets.only(top: 8);
  static const edgeInsetsT12 = EdgeInsets.only(top: 12);
  static const edgeInsetsT24 = EdgeInsets.only(top: 24);

  static const edgeInsetsB4 = EdgeInsets.only(bottom: 4);
  static const edgeInsetsB8 = EdgeInsets.only(bottom: 8);
  static const edgeInsetsB12 = EdgeInsets.only(bottom: 12);
  static const edgeInsetsB24 = EdgeInsets.only(bottom: 24);

  static BorderRadius radius4 = BorderRadius.circular(4);
  static BorderRadius radius8 = BorderRadius.circular(8);
  static BorderRadius radius12 = BorderRadius.circular(12);
  static BorderRadius radius16 = BorderRadius.circular(16);
  static BorderRadius radius24 = BorderRadius.circular(24);
  static BorderRadius radius32 = BorderRadius.circular(32);
  static BorderRadius radius48 = BorderRadius.circular(48);

  // /// 顶部状态栏的高度
  // static double get statusBarHeight => MediaQuery.of(context!).padding.top;
  //
  // /// 底部导航条的高度
  // static double get bottomBarHeight =>
  //     MediaQuery.of(context!).padding.bottom;
}
