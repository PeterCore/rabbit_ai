import 'dart:io';
import 'dart:typed_data';

import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';

import '../app/app_style.dart';

class NetWorkCacheImage extends StatefulWidget {
  final String picUrl;
  final double? width;
  final double? height;
  final BoxFit? fit;
  final double borderRadius;
  final bool progress;
  final bool? isVerticalBorder;
  final Uint8List? bytes;
  const NetWorkCacheImage(this.picUrl,
      {this.width,
      this.height,
      this.bytes,
      this.fit = BoxFit.cover,
      this.borderRadius = 0,
      this.progress = false,
      super.key,
      this.isVerticalBorder = false});

  @override
  State<NetWorkCacheImage> createState() => _NetWorkCacheImageState();
}

class _NetWorkCacheImageState extends State<NetWorkCacheImage>
    with SingleTickerProviderStateMixin {
  late AnimationController animationController;

  @override
  void initState() {
    animationController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 400),
    );
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    var picUrl = widget.picUrl;
    if (picUrl.isEmpty) {
      return Container(
        decoration: BoxDecoration(
          color: Colors.grey.withOpacity(.1),
        ),
        child: const Icon(
          Icons.image,
          color: Colors.grey,
          size: 24,
        ),
      );
    }
    return Container(
      width: widget.width,
      height: widget.height,
      decoration: BoxDecoration(
          borderRadius: widget.isVerticalBorder == true
              ? BorderRadius.vertical(top: Radius.circular(widget.borderRadius))
              : BorderRadius.circular(widget.borderRadius)),
      child: buildImage(url: picUrl, bytes: widget.bytes),
    );
  }

  ExtendedImage buildImage({String url = "", Uint8List? bytes}) {
    if (url.contains("http")) {
      return ExtendedImage.network(
        url,
        fit: widget.fit,
        height: widget.height,
        width: widget.width,
        shape: BoxShape.rectangle,
        handleLoadingProgress: widget.progress,
        borderRadius: widget.isVerticalBorder == true
            ? BorderRadius.vertical(top: Radius.circular(widget.borderRadius))
            : BorderRadius.circular(widget.borderRadius),
        loadStateChanged: (e) {
          if (e.extendedImageLoadState == LoadState.loading) {
            animationController.reset();
            final double? progress =
                e.loadingProgress?.expectedTotalBytes != null
                    ? e.loadingProgress!.cumulativeBytesLoaded /
                        e.loadingProgress!.expectedTotalBytes!
                    : null;
            if (widget.progress) {
              return Center(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    CircularProgressIndicator(
                      value: progress,
                    ),
                    AppStyle.vGap4,
                    Text(
                      '${((progress ?? 0.0) * 100).toInt()}%',
                      textAlign: TextAlign.center,
                      style: const TextStyle(fontSize: 12),
                    ),
                  ],
                ),
              );
            }
            return Container(
              decoration: BoxDecoration(
                color: Colors.grey.withOpacity(.1),
              ),
              child: const Icon(
                Icons.image,
                color: Colors.grey,
                size: 24,
              ),
            );
          }
          if (e.extendedImageLoadState == LoadState.failed) {
            animationController.reset();
            return Container(
              decoration: BoxDecoration(
                color: Colors.grey.withOpacity(.1),
              ),
              child: const Icon(
                Icons.broken_image,
                color: Colors.grey,
                size: 24,
              ),
            );
          }
          if (e.extendedImageLoadState == LoadState.completed) {
            if (e.wasSynchronouslyLoaded) {
              return e.completedWidget;
            }
            animationController.forward();

            return FadeTransition(
              opacity: animationController,
              child: e.completedWidget,
            );
          }
          return null;
        },
      );
    } else if (url.contains("assets")) {
      return ExtendedImage.asset(
        url,
        fit: widget.fit,
        height: widget.height,
        width: widget.width,
        shape: BoxShape.rectangle,
        borderRadius: widget.isVerticalBorder == true
            ? BorderRadius.vertical(top: Radius.circular(widget.borderRadius))
            : BorderRadius.circular(widget.borderRadius),
        loadStateChanged: (e) {
          if (e.extendedImageLoadState == LoadState.loading) {
            animationController.reset();
            final double? progress =
                e.loadingProgress?.expectedTotalBytes != null
                    ? e.loadingProgress!.cumulativeBytesLoaded /
                        e.loadingProgress!.expectedTotalBytes!
                    : null;
            if (widget.progress) {
              return Center(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    CircularProgressIndicator(
                      value: progress,
                    ),
                    AppStyle.vGap4,
                    Text(
                      '${((progress ?? 0.0) * 100).toInt()}%',
                      textAlign: TextAlign.center,
                      style: const TextStyle(fontSize: 12),
                    ),
                  ],
                ),
              );
            }
            return Container(
              decoration: BoxDecoration(
                color: Colors.grey.withOpacity(.1),
              ),
              child: const Icon(
                Icons.image,
                color: Colors.grey,
                size: 24,
              ),
            );
          }
          if (e.extendedImageLoadState == LoadState.failed) {
            animationController.reset();
            return Container(
              decoration: BoxDecoration(
                color: Colors.grey.withOpacity(.1),
              ),
              child: const Icon(
                Icons.broken_image,
                color: Colors.grey,
                size: 24,
              ),
            );
          }
          if (e.extendedImageLoadState == LoadState.completed) {
            if (e.wasSynchronouslyLoaded) {
              return e.completedWidget;
            }
            animationController.forward();

            return FadeTransition(
              opacity: animationController,
              child: e.completedWidget,
            );
          }
          return null;
        },
      );
    } else if (bytes != null) {
      return ExtendedImage.memory(
        bytes,
        fit: widget.fit,
        height: widget.height,
        width: widget.width,
        shape: BoxShape.rectangle,
        borderRadius: widget.isVerticalBorder == true
            ? BorderRadius.vertical(top: Radius.circular(widget.borderRadius))
            : BorderRadius.circular(widget.borderRadius),
        loadStateChanged: (e) {
          if (e.extendedImageLoadState == LoadState.loading) {
            animationController.reset();
            final double? progress =
                e.loadingProgress?.expectedTotalBytes != null
                    ? e.loadingProgress!.cumulativeBytesLoaded /
                        e.loadingProgress!.expectedTotalBytes!
                    : null;
            if (widget.progress) {
              return Center(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    CircularProgressIndicator(
                      value: progress,
                    ),
                    AppStyle.vGap4,
                    Text(
                      '${((progress ?? 0.0) * 100).toInt()}%',
                      textAlign: TextAlign.center,
                      style: const TextStyle(fontSize: 12),
                    ),
                  ],
                ),
              );
            }
            return Container(
              decoration: BoxDecoration(
                color: Colors.grey.withOpacity(.1),
              ),
              child: const Icon(
                Icons.image,
                color: Colors.grey,
                size: 24,
              ),
            );
          }
          if (e.extendedImageLoadState == LoadState.failed) {
            animationController.reset();
            return Container(
              decoration: BoxDecoration(
                color: Colors.grey.withOpacity(.1),
              ),
              child: const Icon(
                Icons.broken_image,
                color: Colors.grey,
                size: 24,
              ),
            );
          }
          if (e.extendedImageLoadState == LoadState.completed) {
            if (e.wasSynchronouslyLoaded) {
              return e.completedWidget;
            }
            animationController.forward();

            return FadeTransition(
              opacity: animationController,
              child: e.completedWidget,
            );
          }
          return null;
        },
      );
    }
    return ExtendedImage.file(
      File(url),
      fit: widget.fit,
      height: widget.height,
      width: widget.width,
      shape: BoxShape.rectangle,
      borderRadius: widget.isVerticalBorder == true
          ? BorderRadius.vertical(top: Radius.circular(widget.borderRadius))
          : BorderRadius.circular(widget.borderRadius),
      loadStateChanged: (e) {
        if (e.extendedImageLoadState == LoadState.loading) {
          animationController.reset();
          final double? progress = e.loadingProgress?.expectedTotalBytes != null
              ? e.loadingProgress!.cumulativeBytesLoaded /
                  e.loadingProgress!.expectedTotalBytes!
              : null;
          if (widget.progress) {
            return Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  CircularProgressIndicator(
                    value: progress,
                  ),
                  AppStyle.vGap4,
                  Text(
                    '${((progress ?? 0.0) * 100).toInt()}%',
                    textAlign: TextAlign.center,
                    style: const TextStyle(fontSize: 12),
                  ),
                ],
              ),
            );
          }
          return Container(
            decoration: BoxDecoration(
              color: Colors.grey.withOpacity(.1),
            ),
            child: const Icon(
              Icons.image,
              color: Colors.grey,
              size: 24,
            ),
          );
        }
        if (e.extendedImageLoadState == LoadState.failed) {
          animationController.reset();
          return Container(
            decoration: BoxDecoration(
              color: Colors.grey.withOpacity(.1),
            ),
            child: const Icon(
              Icons.broken_image,
              color: Colors.grey,
              size: 24,
            ),
          );
        }
        if (e.extendedImageLoadState == LoadState.completed) {
          if (e.wasSynchronouslyLoaded) {
            return e.completedWidget;
          }
          animationController.forward();

          return FadeTransition(
            opacity: animationController,
            child: e.completedWidget,
          );
        }
        return null;
      },
    );
  }

  @override
  void dispose() {
    animationController.dispose();
    super.dispose();
  }
}
