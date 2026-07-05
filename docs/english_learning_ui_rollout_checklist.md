# English Learning UI Rollout Checklist

## Scope
- Auth pages visual refresh: login/register
- Learning pages style unification: home/typing/video-detail/profile
- Home display switch chain: backend model -> web admin form -> uni filtering

## Backend Checks
- Database migration includes `show_home` in `english_learning_video_categories`.
- Database migration includes `show_home` in `english_learning_video_series`.
- Category list API supports query parameter `showHome`.
- Series list API supports query parameter `showHome`.
- `showHome=true` returns only records intended for home display.

## Web Admin Checks
- Video category table shows a "首页展示" status column.
- Video series table shows a "首页展示" status column.
- Category edit/create supports toggling "首页展示".
- Series edit/create supports toggling "首页展示".
- Save and reopen confirms persisted value.

## Uni Home Checks
- Home shows only: app name area + banner + video area.
- Check-in block is hidden.
- Category tabs can switch and refetch data.
- Series list uses 2-column cards.
- Pagination page size is 10.
- Reach-bottom triggers next page until no more data.
- Home request sends `showHome=true`.

## Audio/Typing Regression Checks
- Typing favorites can collect/uncollect normally.
- Typing sentence audio and word audio still playable.
- Audio debug panel still records success/fail/skip entries.
- Settings popup bottom is not blocked by tab area.

## Profile/Language Checks
- Profile language switch uses the same lang-switch component as auth pages.
- Language switching updates visible text without page crash.

## Visual Consistency Checks
- Login/register/profile/typing/video-detail use one visual family.
- Card radius, border style, and CTA gradient are consistent.
- Contrast remains readable in bright environments.
- i18n text wrapping does not break layout in EN/MN/ZH.

## Smoke Test Matrix
- H5 Chrome latest
- Android WebView (uni release build)
- iOS Safari/WebView

## Rollback Notes
- If regression appears only in UI: rollback `uni/src/pages/user/*` and `uni/src/pages/learning/*` style changes.
- If regression appears in home data filtering: rollback `showHome` query usage in uni home and web admin switch fields.
- If migration issue appears: temporarily avoid passing `showHome` in client queries until schema update is complete.
