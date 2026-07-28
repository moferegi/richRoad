<template>
  <div>
    <div class="gva-search-box">
      <el-alert
        title="英语学习运营台（日记端）：支持日记分类、日记管理、字幕解析和用户资源授权。"
        type="info"
        :closable="false"
        show-icon
      />
      <div class="lang-switch-row">
        <span>展示语言：</span>
        <el-select v-model="displayLang" style="width: 140px">
          <el-option
            v-for="lang in displayLanguageOptions"
            :key="lang.value"
            :label="lang.label"
            :value="lang.value"
          />
        </el-select>
      </div>
    </div>

    <div class="gva-table-box">
      <el-tabs v-model="activeTab">
        <!-- 日记分类 -->
        <el-tab-pane label="日记分类" name="diaryCategory">
          <div class="toolbar-row">
            <span class="toolbar-tip">用于维护日记大类（如生活日记、旅行日记等）</span>
            <el-button type="primary" icon="plus" @click="openDiaryCategoryDialog()">新增日记分类</el-button>
          </div>

          <el-table :data="diaryCategoryTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column label="分类名称" min-width="240">
              <template #default="scope">
                {{ formatI18nText(scope.row.name) }}
              </template>
            </el-table-column>
            <el-table-column prop="storageKey" label="存储目录键" min-width="180" show-overflow-tooltip>
              <template #default="scope">
                {{ scope.row.storageKey || '-' }}
              </template>
            </el-table-column>
            <el-table-column label="首页展示" width="120">
              <template #default="scope">
                <el-tag :type="scope.row.showHome === false ? 'info' : 'success'">
                  {{ scope.row.showHome === false ? '隐藏' : '展示' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="sort" label="排序" width="100" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openDiaryCategoryDialog(scope.row)">编辑</el-button>
                <el-button type="danger" link icon="delete" @click="removeDiaryCategory(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="diaryCategoryQuery.page"
              :page-size="diaryCategoryQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="diaryCategoryTotal"
              @current-change="handleDiaryCategoryPageChange"
              @size-change="handleDiaryCategorySizeChange"
            />
          </div>
        </el-tab-pane>

        <!-- 日记管理 -->
        <el-tab-pane label="日记管理" name="diary">
          <div class="toolbar-row">
            <el-form :inline="true" :model="diaryQuery">
              <el-form-item label="所属分类">
                <el-select v-model="diaryQuery.categoryId" clearable placeholder="全部" style="width: 240px">
                  <el-option
                    v-for="item in diaryCategoryOptions"
                    :key="item.ID"
                    :label="formatI18nText(item.name)"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="关键词">
                <el-input v-model="diaryQuery.keyword" clearable placeholder="名称模糊查询" style="width: 220px" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="handleDiarySearch">查询</el-button>
                <el-button icon="refresh" @click="resetDiarySearch">重置</el-button>
              </el-form-item>
            </el-form>
            <el-button type="primary" icon="plus" @click="openDiaryDialog()">新增日记</el-button>
          </div>

          <el-table :data="diaryTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column label="所属分类" min-width="180">
              <template #default="scope">
                {{ diaryCategoryNameMap[scope.row.categoryId] || `#${scope.row.categoryId}` }}
              </template>
            </el-table-column>
            <el-table-column label="日记名称" min-width="220">
              <template #default="scope">
                {{ formatI18nText(scope.row.name) }}
              </template>
            </el-table-column>
            <el-table-column prop="duration" label="时长(秒)" width="100" />
            <el-table-column label="显示翻译" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.showTranslate === false ? 'info' : 'success'">
                  {{ scope.row.showTranslate === false ? '关' : '开' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="显示英文" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.showEnglish === false ? 'info' : 'success'">
                  {{ scope.row.showEnglish === false ? '关' : '开' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="trialPercent" label="试看比例(%)" width="120" />
            <el-table-column prop="sort" label="排序" width="80" />
            <el-table-column label="标签" min-width="160">
              <template #default="scope">
                <el-tag
                  v-for="tag in (scope.row.tags || [])"
                  :key="tag.ID"
                  size="small"
                  :color="tag.color || undefined"
                  style="margin-right: 4px; margin-bottom: 2px;"
                >
                  {{ tag.name }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="280" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openDiaryDialog(scope.row)">编辑</el-button>
                <el-button type="success" link @click="openDiarySubtitleParser(scope.row)">字幕解析</el-button>
                <el-button type="warning" link @click="previewDiarySubtitles(scope.row)">查看字幕</el-button>
                <el-button type="danger" link icon="delete" @click="removeDiary(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="diaryQuery.page"
              :page-size="diaryQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="diaryTotal"
              @current-change="handleDiaryPageChange"
              @size-change="handleDiarySizeChange"
            />
          </div>

          <!-- 字幕解析区 -->
          <div ref="diarySubtitleSectionRef" class="subtitle-box">
            <el-divider content-position="left">日记字幕解析入库</el-divider>
            <el-form label-width="140px">
              <el-form-item label="目标日记">
                <el-select v-model="diarySubtitleForm.diaryId" filterable placeholder="请选择日记" style="width: 480px">
                  <el-option
                    v-for="item in diaryOptions"
                    :key="item.ID"
                    :label="`${diaryCategoryNameMap[item.categoryId] || '-'} / ${formatI18nText(item.name)}`"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="英文字幕文件(必填)">
                <FileUploadWithDir :key="'us-' + diarySubtitleFormKey" v-model="diarySubtitleForm.englishSubtitleUrl" :default-folder="diarySubtitleUploadFolder" accept=".srt,.vtt,.ass" />
              </el-form-item>
              <el-form-item label="英式英文字幕(可选)">
                <FileUploadWithDir :key="'uk-' + diarySubtitleFormKey" v-model="diarySubtitleForm.englishSubtitleUrlUk" :default-folder="diarySubtitleUploadFolder" accept=".srt,.vtt,.ass" />
                <div style="font-size:12px;color:#999;margin-top:4px;">仅时间轴与美式不同时上传，文字内容以美式为准</div>
              </el-form-item>

              <el-form-item
                v-for="lang in subtitleLanguages"
                :key="lang.code"
                :label="`${lang.name || lang.code} 字幕文件`"
              >
                <FileUploadWithDir :key="'trans-' + lang.code + '-' + diarySubtitleFormKey" v-model="diarySubtitleForm.translationSubtitleMap[lang.code]" :default-folder="diarySubtitleUploadFolder" accept=".srt,.vtt,.ass" />
              </el-form-item>

              <el-form-item>
                <el-button type="warning" :loading="diaryKeywordScanning" @click="handleDiaryScanKeywords">
                  {{ diaryKeywordScanning ? '正在扫描...' : '步骤1: 扫描关键词' }}
                </el-button>
                <el-button
                  v-if="diaryKeywordScanDone"
                  type="primary"
                  style="margin-left: 12px"
                  :disabled="diarySelectedKeywordIds.length === 0"
                  @click="handleDiaryConfirmParse"
                >
                  步骤2: 确认解析入库（已选 {{ diarySelectedKeywordIds.length }} 个词）
                </el-button>
              </el-form-item>
            </el-form>

            <!-- 关键词扫描结果 -->
            <div v-if="diaryKeywordScanDone" class="keyword-scan-result">
              <el-divider content-position="left">
                关键词扫描结果（共 {{ diaryKeywordScanResults.length }} 个单词，已选 {{ diarySelectedKeywordIds.length }} 个）
              </el-divider>
              <div class="keyword-toolbar">
                <el-button size="small" @click="diarySelectAllKeywords(true)">全选已匹配</el-button>
                <el-button size="small" @click="diarySelectAllKeywords(false)">取消全选</el-button>
              </div>
              <el-table
                ref="diaryKeywordTableRef"
                :data="diaryKeywordScanResults"
                row-key="word"
                border
                max-height="400"
                @selection-change="handleDiaryKeywordSelectionChange"
              >
                <el-table-column type="selection" width="50" :selectable="(row) => row.matched" />
                <el-table-column prop="word" label="单词" width="180" />
                <el-table-column label="状态" width="120">
                  <template #default="scope">
                    <el-tag v-if="scope.row.matched" type="success" size="small">已匹配(ID:{{ scope.row.wordId }})</el-tag>
                    <el-tag v-else type="info" size="small">未入库</el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="count" label="出现次数" width="100" sortable />
              </el-table>
            </div>
          </div>
        </el-tab-pane>

        <!-- 用户授权 -->
        <el-tab-pane label="用户授权" name="entitlement">
          <div class="toolbar-row">
            <el-form :inline="true" :model="entitlementQuery">
              <el-form-item label="用户ID">
                <el-input-number v-model="entitlementQuery.userId" :min="1" controls-position="right" />
              </el-form-item>
              <el-form-item label="资源类型">
                <el-select v-model="entitlementQuery.resourceType" clearable placeholder="全部" style="width: 220px">
                  <el-option label="日记分类" value="english_category" />
                  <el-option label="日记" value="diary" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="handleEntitlementSearch">查询</el-button>
                <el-button icon="refresh" @click="resetEntitlementSearch">重置</el-button>
              </el-form-item>
            </el-form>
            <el-button type="primary" icon="plus" @click="openEntitlementDialog">新增授权</el-button>
          </div>

          <el-table :data="entitlementTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column prop="userId" label="用户ID" width="100" />
            <el-table-column label="资源类型" width="140">
              <template #default="scope">
                <el-tag>{{ resourceTypeLabel(scope.row.resourceType) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="resourceId" label="资源ID" width="110" />
            <el-table-column prop="source" label="授权来源" width="140" />
            <el-table-column prop="grantedBy" label="授权人" width="100" />
            <el-table-column label="过期时间" min-width="180">
              <template #default="scope">
                {{ scope.row.expireAt || '永久有效' }}
              </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注" min-width="200" show-overflow-tooltip />
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="scope">
                <el-button type="danger" link icon="delete" @click="removeEntitlement(scope.row)">撤销</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="entitlementQuery.page"
              :page-size="entitlementQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="entitlementTotal"
              @current-change="handleEntitlementPageChange"
              @size-change="handleEntitlementSizeChange"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 日记分类弹窗 -->
    <el-dialog v-model="diaryCategoryDialogVisible" :title="diaryCategoryDialogMode === 'create' ? '新增日记分类' : '编辑日记分类'" width="760px" :close-on-click-modal="false" :close-on-press-escape="false">
      <el-form :model="diaryCategoryForm" label-width="120px">
        <el-form-item label="分类名称">
          <MultiLangEditor
            :model="diaryCategoryForm.nameI18n"
            title="日记分类多语言名称"
            input-type="input"
            :rows="2"
            :use-tabs="true"
          />
        </el-form-item>
        <el-form-item label="存储目录键">
          <el-input v-model="diaryCategoryForm.storageKey" placeholder="例如: life-diary" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="diaryCategoryForm.sort" :min="0" />
        </el-form-item>
        <el-form-item label="首页展示">
          <el-switch v-model="diaryCategoryForm.showHome" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="diaryCategoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitDiaryCategory">确定</el-button>
      </template>
    </el-dialog>

    <!-- 日记弹窗 -->
    <el-dialog v-model="diaryDialogVisible" :title="diaryDialogMode === 'create' ? '新增日记' : '编辑日记'" width="860px" :close-on-click-modal="false" :close-on-press-escape="false">
      <el-form :model="diaryForm" label-width="120px">
        <el-form-item label="所属分类">
          <el-select v-model="diaryForm.categoryId" filterable placeholder="请选择分类" style="width: 100%">
            <el-option
              v-for="item in diaryCategoryOptions"
              :key="item.ID"
              :label="formatI18nText(item.name)"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="日记名称">
          <MultiLangEditor
            :model="diaryForm.nameI18n"
            title="日记名称多语言"
            input-type="input"
            :rows="2"
            :use-tabs="true"
          />
        </el-form-item>
        <el-form-item label="美式语音(US)">
          <FileUploadWithDir v-model="diaryForm.audioUs" default-folder="english-learn/diary/audio" accept="audio/*" />
        </el-form-item>
        <el-form-item label="英式语音(UK)">
          <FileUploadWithDir v-model="diaryForm.audioUk" default-folder="english-learn/diary/audio" accept="audio/*" />
        </el-form-item>
        <el-form-item label="音频时长(秒)">
          <template #label>
            音频时长(秒)
            <el-tooltip content="上传美式语音后自动计算时长，也可手动修改">
              <el-icon style="margin-left:4px;vertical-align:middle;"><QuestionFilled /></el-icon>
            </el-tooltip>
          </template>
          <el-input-number v-model="diaryForm.duration" :min="0" :precision="1" />
        </el-form-item>
        <el-form-item label="日记图片(多语言)">
          <FileUploadWithDir v-model="diaryForm.imageI18n" default-folder="english-learn/diary/image" accept="image/*" />
        </el-form-item>
        <el-form-item label="显示翻译字幕">
          <el-switch v-model="diaryForm.showTranslate" />
        </el-form-item>
        <el-form-item label="显示英文字幕">
          <el-switch v-model="diaryForm.showEnglish" />
        </el-form-item>
        <el-form-item label="试看比例(%)">
          <el-input-number v-model="diaryForm.trialPercent" :min="1" :max="100" />
        </el-form-item>
        <el-form-item label="需要会员">
          <el-switch v-model="diaryForm.needVip" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="diaryForm.sort" :min="0" />
        </el-form-item>
        <el-form-item label="日记标签">
          <el-select
            v-model="diaryForm.tagIds"
            multiple
            filterable
            clearable
            placeholder="请选择日记标签"
            style="width: 100%"
          >
            <el-option
              v-for="item in diaryTagOptions"
              :key="item.ID"
              :label="item.name"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="diaryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitDiary">确定</el-button>
      </template>
    </el-dialog>

    <!-- 授权弹窗 -->
    <el-dialog v-model="entitlementDialogVisible" title="新增资源授权" width="760px" :close-on-click-modal="false" :close-on-press-escape="false">
      <el-form :model="entitlementForm" label-width="120px">
        <el-form-item label="用户ID">
          <el-input-number v-model="entitlementForm.userId" :min="1" controls-position="right" />
        </el-form-item>
        <el-form-item label="资源类型">
          <el-select v-model="entitlementForm.resourceType" style="width: 100%">
            <el-option label="日记分类" value="english_category" />
            <el-option label="日记" value="diary" />
          </el-select>
        </el-form-item>
        <el-form-item label="资源ID">
          <el-select
            v-model="entitlementForm.resourceId"
            filterable
            clearable
            placeholder="可从列表选择，或手动输入"
            style="width: 100%"
          >
            <el-option
              v-for="item in resourceOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
          <el-input-number v-model="entitlementForm.resourceId" :min="1" class="mt-2 w-full" controls-position="right" />
        </el-form-item>
        <el-form-item label="过期时间">
          <el-date-picker
            v-model="entitlementForm.expireAt"
            type="datetime"
            placeholder="不填表示永久"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
            style="width: 100%"
            clearable
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="entitlementForm.remark" type="textarea" :rows="3" placeholder="可选" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="entitlementDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEntitlement">确定</el-button>
      </template>
    </el-dialog>

    <!-- 字幕预览弹窗 -->
    <el-dialog
      v-model="diarySentencePreviewVisible"
      :title="`字幕预览 - ${diaryPreviewDiaryName || ''}`"
      width="1100px"
      :before-close="handleDiarySentencePreviewBeforeClose"
    >
      <el-tabs v-model="diarySentencePreviewTab">
        <el-tab-pane label="字幕句子" name="sentences">
          <div class="subtitle-preview-toolbar">
            <span>字幕语言：</span>
            <el-select v-model="diarySentencePreviewLang" style="width: 140px">
              <el-option
                v-for="lang in displayLanguageOptions"
                :key="lang.value"
                :label="lang.label"
                :value="lang.value"
              />
            </el-select>
            <el-tag type="warning" effect="light">已改动 {{ diaryPendingSentenceChanges.length }} 条</el-tag>
          </div>
          <el-table :data="diarySentencePreviewTable" border max-height="500">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column label="开始(s)" width="120">
              <template #default="scope">
                <el-input-number v-model="scope.row.startTime" :min="0" :step="0.1" controls-position="right" @change="diaryPendingSentenceChanges = computeDiaryPendingChanges()" />
              </template>
            </el-table-column>
            <el-table-column label="结束(s)" width="120">
              <template #default="scope">
                <el-input-number v-model="scope.row.endTime" :min="0" :step="0.1" controls-position="right" @change="diaryPendingSentenceChanges = computeDiaryPendingChanges()" />
              </template>
            </el-table-column>
            <el-table-column label="英文句子" min-width="320">
              <template #default="scope">
                <el-input v-model="scope.row.english" type="textarea" :rows="2" @input="diaryPendingSentenceChanges = computeDiaryPendingChanges()" />
              </template>
            </el-table-column>
            <el-table-column label="翻译" min-width="320">
              <template #default="scope">
                <el-input v-model="scope.row.translateText" type="textarea" :rows="2" @input="diarySyncSentenceTranslate(scope.row)" />
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="重点单词" name="keywords">
          <div class="subtitle-preview-toolbar">
            <span>字幕单词（共 {{ diaryPreviewKeywordList.length }} 个，已匹配词库 {{ diaryPreviewKeywordList.filter(k => k.matched).length }} 个）</span>
            <el-button size="small" type="primary" @click="loadDiaryPreviewKeywords">刷新列表</el-button>
          </div>
          <el-table
            ref="diaryPreviewKeywordTableRef"
            :data="diaryPreviewKeywordList"
            row-key="word"
            border
            max-height="420"
            @selection-change="handleDiaryPreviewKeywordSelectionChange"
          >
            <el-table-column type="selection" width="50" :selectable="(row) => row.matched" />
            <el-table-column prop="word" label="单词" width="180" />
            <el-table-column label="状态" width="140">
              <template #default="scope">
                <el-tag v-if="scope.row.matched" type="success" size="small">已匹配(ID:{{ scope.row.wordId }})</el-tag>
                <el-tag v-else type="info" size="small">未入库</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="count" label="出现次数" width="100" sortable />
          </el-table>
          <div style="margin-top: 12px">
            <el-button size="small" @click="diarySelectAllPreviewKeywords(true)">全选已匹配</el-button>
            <el-button size="small" @click="diarySelectAllPreviewKeywords(false)">取消全选</el-button>
            <el-button
              type="warning"
              size="small"
              style="margin-left: 16px"
              :loading="diaryPreviewRehighlighting"
              :disabled="diaryPreviewSelectedKeywordIds.length === 0"
              @click="handleDiaryRehighlightPreview"
            >
              更新高亮（已选 {{ diaryPreviewSelectedKeywordIds.length }} 个词）
            </el-button>
          </div>
        </el-tab-pane>

        <el-tab-pane label="全部单词" name="allWords">
          <div class="subtitle-preview-toolbar">
            <span>当前字幕中的所有单词（共 {{ allDiarySubtitleWords.length }} 个），勾选后点击"更新高亮"即可添加为重点单词</span>
          </div>
          <el-table
            ref="diaryAllWordsTableRef"
            :data="allDiarySubtitleWords"
            row-key="word"
            border
            max-height="420"
            @selection-change="handleDiaryAllWordsSelectionChange"
          >
            <el-table-column type="selection" width="50" />
            <el-table-column prop="word" label="单词" width="180" />
            <el-table-column label="已在词库" width="140">
              <template #default="scope">
                <el-tag v-if="scope.row.matched" type="success" size="small">是(ID:{{ scope.row.wordId }})</el-tag>
                <el-tag v-else type="info" size="small">否</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="count" label="出现次数" width="100" sortable />
          </el-table>
          <div style="margin-top: 12px">
            <el-button size="small" @click="diarySelectAllAllWords(true)">全选已匹配</el-button>
            <el-button size="small" @click="diarySelectAllAllWords(false)">取消全选</el-button>
            <el-button
              type="warning"
              size="small"
              style="margin-left: 16px"
              :loading="diaryAllWordsRehighlighting"
              :disabled="diaryAllWordsSelectedIds.length === 0"
              @click="handleDiaryRehighlightAllWords"
            >
              更新高亮（已选 {{ diaryAllWordsSelectedIds.length }} 个词）
            </el-button>
          </div>
        </el-tab-pane>
      </el-tabs>

      <template #footer>
        <el-button :disabled="diaryPendingSentenceChanges.length === 0" @click="diaryResetSentencePreviewChanges">重置改动</el-button>
        <el-button type="primary" :loading="diarySentenceSaving" :disabled="diaryPendingSentenceChanges.length === 0" @click="diarySaveSentencePreview">保存修改</el-button>
        <el-button @click="diaryCloseSentencePreview">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, nextTick, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { QuestionFilled } from '@element-plus/icons-vue'
import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'
import FileUploadWithDir from '@/components/FileUploadWithDir/index.vue'
import { getLanguageList } from '@/api/client/language'
import { getDiaryTagList } from '@/api/client/diaryTag'
import {
  createDiaryCategory,
  createDiary,
  updateDiaryCategory,
  updateDiary,
  deleteDiaryCategory,
  deleteDiary,
  findDiary,
  getDiaryCategoryList,
  getDiaryList,
  getDiarySentenceList,
  getDiarySentenceListAll,
  parseDiarySubtitle,
  scanDiaryKeywords,
  parseDiarySubtitleFiles,
  getEntitlementList,
  grantEntitlement,
  revokeEntitlement,
  rehighlightDiarySentences,
  updateDiarySentenceList,
  getDiaryKeywords,
  batchCreateWords
} from '../api/english'

defineOptions({
  name: 'EnglishLearningDiary'
})

const activeTab = ref('diaryCategory')
const displayLang = ref('zh')

const formatI18nText = (raw) => {
  if (!raw) return ''
  try {
    const parsed = typeof raw === 'string' ? JSON.parse(raw) : raw
    if (typeof parsed === 'object' && parsed !== null) {
      const lang = String(displayLang.value || 'zh')
      return parsed[lang] || parsed.zh || parsed.en || Object.values(parsed)[0] || ''
    }
    return String(raw)
  } catch (e) {
    return String(raw)
  }
}

const stringifyI18nObject = (i18nObject) => {
  const source = i18nObject && typeof i18nObject === 'object' ? i18nObject : {}
  const cleaned = {}
  for (const [key, value] of Object.entries(source)) {
    const lang = String(key || '').trim()
    if (!lang) continue
    const text = String(value ?? '').trim()
    if (text) cleaned[lang] = text
  }
  return JSON.stringify(Object.keys(cleaned).length > 0 ? cleaned : { zh: '' })
}

// 从日记图片JSON字段中提取URL
const extractImageUrl = (imageI18n) => {
  if (!imageI18n) return ''
  try {
    const parsed = typeof imageI18n === 'object' ? imageI18n : JSON.parse(imageI18n || '{}')
    return parsed.zh || parsed.en || Object.values(parsed)[0] || ''
  } catch (e) {
    return String(imageI18n || '')
  }
}

const managedLanguages = ref([])
const displayLanguageOptions = computed(() => {
  if (managedLanguages.value.length === 0) {
    return [
      { value: 'zh', label: '中文' },
      { value: 'en', label: 'English' }
    ]
  }
  return managedLanguages.value.map((lang) => ({
    value: lang.code,
    label: lang.name || lang.nativeName || lang.code
  }))
})

// ===== 日记分类 =====
const diaryCategoryQuery = ref({ page: 1, pageSize: 10 })
const diaryCategoryTable = ref([])
const diaryCategoryTotal = ref(0)
const diaryCategoryOptions = ref([])

const diaryCategoryNameMap = computed(() => {
  const map = {}
  for (const item of diaryCategoryOptions.value) {
    map[item.ID] = formatI18nText(item.name)
  }
  return map
})

const loadDiaryCategoryList = async () => {
  const res = await getDiaryCategoryList({ page: diaryCategoryQuery.value.page, pageSize: diaryCategoryQuery.value.pageSize })
  if (res.code === 0 && res.data) {
    diaryCategoryTable.value = res.data.list || []
    diaryCategoryTotal.value = res.data.total || 0
  }
}

const loadDiaryCategoryOptions = async () => {
  const res = await getDiaryCategoryList({ page: 1, pageSize: 200 })
  if (res.code === 0 && res.data) {
    diaryCategoryOptions.value = res.data.list || []
  }
}

const handleDiaryCategoryPageChange = (page) => {
  diaryCategoryQuery.value.page = page
  loadDiaryCategoryList()
}
const handleDiaryCategorySizeChange = (size) => {
  diaryCategoryQuery.value.pageSize = size
  diaryCategoryQuery.value.page = 1
  loadDiaryCategoryList()
}

const diaryCategoryDialogVisible = ref(false)
const diaryCategoryDialogMode = ref('create')
const diaryCategoryForm = ref({
  ID: 0,
  nameI18n: { zh: '' },
  storageKey: '',
  showHome: true,
  sort: 0
})

const openDiaryCategoryDialog = (row) => {
  if (row) {
    diaryCategoryDialogMode.value = 'update'
    diaryCategoryForm.value = {
      ID: row.ID,
      nameI18n: JSON.parse(row.name || '{"zh":""}'),
      storageKey: row.storageKey || '',
      showHome: row.showHome !== false,
      sort: row.sort || 0
    }
  } else {
    diaryCategoryDialogMode.value = 'create'
    diaryCategoryForm.value = { ID: 0, nameI18n: { zh: '' }, storageKey: '', showHome: true, sort: 0 }
  }
  diaryCategoryDialogVisible.value = true
}

const submitDiaryCategory = async () => {
  const data = {
    name: stringifyI18nObject(diaryCategoryForm.value.nameI18n),
    storageKey: diaryCategoryForm.value.storageKey,
    showHome: diaryCategoryForm.value.showHome,
    sort: diaryCategoryForm.value.sort
  }
  let res
  if (diaryCategoryDialogMode.value === 'create') {
    res = await createDiaryCategory(data)
  } else {
    data.ID = diaryCategoryForm.value.ID
    res = await updateDiaryCategory(data)
  }
  if (res.code === 0) {
    ElMessage.success(diaryCategoryDialogMode.value === 'create' ? '创建成功' : '更新成功')
    diaryCategoryDialogVisible.value = false
    loadDiaryCategoryList()
    loadDiaryCategoryOptions()
  }
}

const removeDiaryCategory = (row) => {
  ElMessageBox.confirm('确定要删除该日记分类吗？', '提示', { type: 'warning' }).then(async () => {
    const res = await deleteDiaryCategory({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      loadDiaryCategoryList()
      loadDiaryCategoryOptions()
    }
  })
}

// ===== 日记 =====
const diaryQuery = ref({ page: 1, pageSize: 10, categoryId: undefined, keyword: '' })
const diaryTable = ref([])
const diaryTotal = ref(0)
const diaryOptions = ref([])
const diaryTagOptions = ref([])

const loadDiaryList = async () => {
  const res = await getDiaryList({ page: diaryQuery.value.page, pageSize: diaryQuery.value.pageSize, categoryId: diaryQuery.value.categoryId, keyword: diaryQuery.value.keyword })
  if (res.code === 0 && res.data) {
    diaryTable.value = res.data.list || []
    diaryTotal.value = res.data.total || 0
  }
}

const loadDiaryOptions = async () => {
  const res = await getDiaryList({ page: 1, pageSize: 500 })
  if (res.code === 0 && res.data) {
    diaryOptions.value = res.data.list || []
  }
}

const loadDiaryTagOptions = async () => {
  const res = await getDiaryTagList({ page: 1, pageSize: 200 })
  if (res.code === 0 && res.data) {
    diaryTagOptions.value = res.data.list || []
  }
}

const handleDiarySearch = () => { diaryQuery.value.page = 1; loadDiaryList() }
const resetDiarySearch = () => {
  diaryQuery.value = { page: 1, pageSize: 10, categoryId: undefined, keyword: '' }
  loadDiaryList()
}
const handleDiaryPageChange = (page) => { diaryQuery.value.page = page; loadDiaryList() }
const handleDiarySizeChange = (size) => { diaryQuery.value.pageSize = size; diaryQuery.value.page = 1; loadDiaryList() }

const diaryDialogVisible = ref(false)
const diaryDialogMode = ref('create')
const diaryForm = ref({
  ID: 0,
  categoryId: undefined,
  nameI18n: { zh: '' },
  audioUs: '',
  audioUk: '',
  duration: 0,
  imageI18n: { zh: '' },
  showTranslate: true,
  showEnglish: true,
  trialPercent: 8,
  needVip: false,
  englishSubtitleUrl: '',
  englishSubtitleUrlUk: '',
  sort: 0,
  tagIds: []
})

// 从音频URL获取时长
const getAudioDuration = (url) => {
  return new Promise((resolve) => {
    if (!url) {
      resolve(0)
      return
    }
    const audio = new Audio()
    audio.preload = 'metadata'
    audio.crossOrigin = 'anonymous'
    let resolved = false
    const cleanup = () => {
      if (resolved) return
      resolved = true
      audio.removeEventListener('loadedmetadata', onLoaded)
      audio.removeEventListener('error', onError)
      audio.src = ''
    }
    const onLoaded = () => {
      cleanup()
      resolve(audio.duration && isFinite(audio.duration) ? audio.duration : 0)
    }
    const onError = () => {
      cleanup()
      console.warn('获取音频时长失败，可能是CORS限制:', url)
      resolve(0)
    }
    audio.addEventListener('loadedmetadata', onLoaded)
    audio.addEventListener('error', onError)
    setTimeout(() => cleanup(), 10000)
    audio.src = url
  })
}

// 监听美式语音URL变化，自动填充时长
watch(() => diaryForm.value.audioUs, async (newUrl, oldUrl) => {
  if (!newUrl || newUrl === oldUrl) return
  const duration = await getAudioDuration(newUrl)
  if (duration > 0) {
    diaryForm.value.duration = Math.round(duration * 10) / 10
  }
})

const openDiaryDialog = async (row) => {
  if (row) {
    diaryDialogMode.value = 'update'
    diaryForm.value = {
      ID: row.ID,
      categoryId: row.categoryId,
      nameI18n: JSON.parse(row.name || '{"zh":""}'),
      audioUs: row.audioUs || '',
      audioUk: row.audioUk || '',
      duration: row.duration || 0,
      imageI18n: extractImageUrl(row.imageI18n),
      showTranslate: row.showTranslate !== false,
      showEnglish: row.showEnglish !== false,
      trialPercent: row.trialPercent || 8,
      needVip: row.needVip || false,
      englishSubtitleUrl: row.englishSubtitleUrl || '',
      englishSubtitleUrlUk: row.englishSubtitleUrlUk || '',
      sort: row.sort || 0,
      tagIds: []
    }
    // Call findDiary API to load tags
    try {
      const detailRes = await findDiary({ ID: row.ID })
      if (detailRes.code === 0 && detailRes.data?.tags) {
        diaryForm.value.tagIds = detailRes.data.tags.map(t => t.ID)
      }
    } catch (e) {
      // ignore tag load failure
    }
  } else {
    diaryDialogMode.value = 'create'
    diaryForm.value = {
      ID: 0, categoryId: undefined, nameI18n: { zh: '' }, audioUs: '', audioUk: '',
      duration: 0, imageI18n: '', showTranslate: true, showEnglish: true,
      trialPercent: 8, needVip: false, englishSubtitleUrl: '', englishSubtitleUrlUk: '',
      sort: 0, tagIds: []
    }
  }
  diaryDialogVisible.value = true
}

const submitDiary = async () => {
  // 日记图片URL存入JSON格式（兼容uni端多语言读取）
  const imageUrl = diaryForm.value.imageI18n || ''
  const imageI18nJson = imageUrl ? JSON.stringify({ zh: imageUrl }) : '{}'
  const data = {
    categoryId: diaryForm.value.categoryId,
    name: stringifyI18nObject(diaryForm.value.nameI18n),
    audioUs: diaryForm.value.audioUs,
    audioUk: diaryForm.value.audioUk,
    duration: diaryForm.value.duration,
    imageI18n: imageI18nJson,
    showTranslate: diaryForm.value.showTranslate,
    showEnglish: diaryForm.value.showEnglish,
    trialPercent: diaryForm.value.trialPercent,
    needVip: diaryForm.value.needVip,
    englishSubtitleUrl: diaryForm.value.englishSubtitleUrl,
    englishSubtitleUrlUk: diaryForm.value.englishSubtitleUrlUk,
    sort: diaryForm.value.sort,
    tagIds: diaryForm.value.tagIds
  }
  let res
  if (diaryDialogMode.value === 'create') {
    res = await createDiary(data)
  } else {
    data.ID = diaryForm.value.ID
    res = await updateDiary(data)
  }
  if (res.code === 0) {
    ElMessage.success(diaryDialogMode.value === 'create' ? '创建成功' : '更新成功')
    diaryDialogVisible.value = false
    loadDiaryList()
    loadDiaryOptions()
  }
}

const removeDiary = (row) => {
  ElMessageBox.confirm('确定要删除该日记吗？', '提示', { type: 'warning' }).then(async () => {
    const res = await deleteDiary({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      loadDiaryList()
      loadDiaryOptions()
    }
  })
}

// 字幕预览
const diarySentencePreviewVisible = ref(false)
const diarySentencePreviewTab = ref('sentences')
const diarySentencePreviewLang = ref('zh')
const diarySentencePreviewTable = ref([])
const diarySentencePreviewOriginalMap = ref({})
const diaryPreviewDiaryId = ref(0)
const diaryPreviewDiaryName = ref('')
const diaryPendingSentenceChanges = ref([])
const diarySentenceSaving = ref(false)
const diaryPreviewKeywordList = ref([])
const diaryPreviewSelectedKeywordIds = ref([])
const diaryPreviewRehighlighting = ref(false)
const diaryPreviewKeywordTableRef = ref(null)

// 全部单词 tab - 状态
const diaryAllWordsTableRef = ref(null)
const diaryAllWordsSelectedIds = ref([])
const diaryAllWordsRehighlighting = ref(false)
const diarySkipAutoSelectAllWords = ref(false)

// 全部单词：从当前字幕句子中提取所有单词并统计词频
const allDiarySubtitleWords = computed(() => {
  const wordCounter = new Map()
  const wordRegex = /[a-zA-Z]+/g
  const highlightedWordMap = new Map() // 在词库中的词 -> wordId
  const wTagRegex = /<w id="(\d+)">([^<]+)<\/w>/g
  const sentenceHighlighted = new Set() // 句子中实际被 <w> 标签高亮的词
  for (const kw of diaryPreviewKeywordList.value) {
    if (kw.matched) {
      highlightedWordMap.set(kw.word.toLowerCase(), kw.wordId)
    }
  }
  for (const row of diarySentencePreviewTable.value) {
    const rawText = String(row?.english || '')
    // 提取 <w> 标签中的词，记录为实际高亮
    let m
    while ((m = wTagRegex.exec(rawText)) !== null) {
      sentenceHighlighted.add(m[2].toLowerCase())
    }
    wTagRegex.lastIndex = 0
    // 去标签后统计所有单词出现次数
    const text = rawText.replace(/<[^>]+>/g, '')
    const matches = text.matchAll(wordRegex)
    for (const match of matches) {
      const word = match[0].toLowerCase()
      wordCounter.set(word, (wordCounter.get(word) || 0) + 1)
    }
  }
  const entries = Array.from(wordCounter.entries())
    .map(([word, count]) => ({ word, count }))
    .sort((a, b) => b.count - a.count)
  return entries.map(({ word, count }) => ({
    word,
    count,
    matched: highlightedWordMap.has(word),
    wordId: highlightedWordMap.get(word) || 0,
    isHighlighted: sentenceHighlighted.has(word)
  }))
})

const parseI18nObject = (value) => {
  if (!value) return {}
  if (typeof value === 'object') return value
  try {
    return JSON.parse(value)
  } catch (e) {
    return {}
  }
}

const buildDiarySentencePreviewSnapshotMap = (list) => {
  const map = {}
  for (const item of list) {
    const rowId = Number(item?.id || item?.ID || 0)
    if (!rowId) continue
    const translateObj = parseI18nObject(item?.translate)
    map[rowId] = {
      startTime: Number(item?.startTime || 0),
      endTime: Number(item?.endTime || 0),
      english: String(item?.english || '').trim(),
      translateObj: JSON.stringify(translateObj)
    }
  }
  return map
}

const diarySyncSentenceTranslate = (row) => {
  if (!row || typeof row !== 'object') return
  const lang = String(diarySentencePreviewLang.value || 'zh')
  const next = { ...(row.translateObj || {}) }
  next[lang] = String(row.translateText || '')
  row.translateObj = next
  row.translate = JSON.stringify(next)
  diaryPendingSentenceChanges.value = computeDiaryPendingChanges()
}

const diaryRefreshSentenceTranslateByLang = () => {
  const lang = String(diarySentencePreviewLang.value || 'zh')
  diarySentencePreviewTable.value = diarySentencePreviewTable.value.map((item) => {
    const translateObj = parseI18nObject(item?.translateObj || item?.translate)
    return {
      ...item,
      translateObj,
      translateText: translateObj[lang] || translateObj.zh || translateObj.en || translateObj.mn || ''
    }
  })
}

// 预览语言切换时刷新翻译
watch(diarySentencePreviewLang, () => {
  diaryRefreshSentenceTranslateByLang()
})

const loadDiaryPreviewKeywords = async () => {
  if (!diaryPreviewDiaryId.value) return
  // 调用后端接口获取完整词库列表（含已高亮状态和出现次数）
  const res = await getDiaryKeywords({ diaryId: diaryPreviewDiaryId.value })
  if (res.code !== 0) return
  diaryPreviewKeywordList.value = (res.data || []).map(item => ({
    word: item.word,
    wordId: item.wordId,
    matched: !!item.matched,
    count: item.count || 0
  }))
  // 自动激活已匹配词库的单词（选中表格行）
  nextTick(() => {
    if (diaryPreviewKeywordTableRef.value) {
      const matchedRows = diaryPreviewKeywordList.value.filter(k => k.matched)
      matchedRows.forEach(row => {
        diaryPreviewKeywordTableRef.value.toggleRowSelection(row, true)
      })
      diaryPreviewSelectedKeywordIds.value = matchedRows.map(r => r.wordId).filter(Boolean)
    }
  })
}

const handleDiaryPreviewKeywordSelectionChange = (rows) => {
  diaryPreviewSelectedKeywordIds.value = rows.map(r => r.wordId).filter(Boolean)
}

const diarySelectAllPreviewKeywords = (select) => {
  if (!diaryPreviewKeywordTableRef.value) return
  if (select) {
    const matched = diaryPreviewKeywordList.value.filter(k => k.matched)
    matched.forEach(row => diaryPreviewKeywordTableRef.value.toggleRowSelection(row, true))
  } else {
    diaryPreviewKeywordTableRef.value.clearSelection()
  }
}

const handleDiaryRehighlightPreview = async () => {
  if (!diaryPreviewDiaryId.value || diaryPreviewSelectedKeywordIds.value.length === 0) return
  diaryPreviewRehighlighting.value = true
  try {
    const res = await rehighlightDiarySentences({
      diaryId: diaryPreviewDiaryId.value,
      keywordIds: diaryPreviewSelectedKeywordIds.value
    })
    if (res.code === 0) {
      ElMessage.success('重新高亮成功')
      // 刷新字幕句子表
      const sentRes = await getDiarySentenceListAll({ diaryId: diaryPreviewDiaryId.value })
      if (sentRes.code === 0) {
        diarySentencePreviewTable.value = (Array.isArray(sentRes.data) ? sentRes.data : []).map((item) => {
          const translateObj = parseI18nObject(item?.translate)
          const lang = String(diarySentencePreviewLang.value || 'zh')
          const normalizedTranslate = typeof item?.translate === 'string'
            ? item.translate
            : JSON.stringify(translateObj)
          return {
            ...item,
            translate: normalizedTranslate,
            translateObj,
            translateText: translateObj[lang] || translateObj.zh || translateObj.en || translateObj.mn || ''
          }
        })
        // 重置快照和改动计数
        diarySentencePreviewOriginalMap.value = buildDiarySentencePreviewSnapshotMap(diarySentencePreviewTable.value)
        diaryPendingSentenceChanges.value = []
      }
      // 刷新关键词列表
      diarySkipAutoSelectAllWords.value = true
      await loadDiaryPreviewKeywords()
      diarySkipAutoSelectAllWords.value = false
    } else {
      ElMessage.error(res.msg || '重新高亮失败')
    }
  } catch (e) {
    ElMessage.error('重新高亮失败')
  }
  diaryPreviewRehighlighting.value = false
}

// 全部单词 tab - 选择管理
const handleDiaryAllWordsSelectionChange = (rows) => {
  diaryAllWordsSelectedIds.value = rows.map(r => r.wordId).filter(id => id > 0)
}

const diarySelectAllAllWords = (selectAll) => {
  if (!diaryAllWordsTableRef.value) return
  if (selectAll) {
    const matchedRows = allDiarySubtitleWords.value.filter(k => k.matched)
    matchedRows.forEach(row => diaryAllWordsTableRef.value.toggleRowSelection(row, true))
  } else {
    diaryAllWordsTableRef.value.clearSelection()
  }
}

const handleDiaryRehighlightAllWords = async () => {
  if (!diaryPreviewDiaryId.value) {
    ElMessage.warning('日记ID无效')
    return
  }
  if (diaryAllWordsSelectedIds.value.length === 0) {
    ElMessage.warning('请至少选择一个已入库的单词')
    return
  }

  diaryAllWordsRehighlighting.value = true

  // 收集所有选中行，找出未入库的单词
  const selectedRows = diaryAllWordsTableRef.value?.getSelectionRows?.() || []
  const unmatchedWords = selectedRows
    .filter(r => !r.matched)
    .map(r => r.word)

  let newIds = []
  if (unmatchedWords.length > 0) {
    const batchRes = await batchCreateWords({ words: unmatchedWords })
    if (batchRes.code === 0 && Array.isArray(batchRes.data)) {
      newIds = batchRes.data.map(item => item.wordId)
      ElMessage.success(`已自动入库 ${newIds.length} 个新单词`)
    }
  }

  const allIds = [...new Set([...diaryAllWordsSelectedIds.value, ...newIds])]
  const res = await rehighlightDiarySentences({
    diaryId: diaryPreviewDiaryId.value,
    keywordIds: allIds
  })
  diaryAllWordsRehighlighting.value = false

  if (res.code !== 0) return
  ElMessage.success(`字幕重新高亮成功，已应用 ${allIds.length} 个重点单词`)

  // 刷新句子列表
  const sentRes = await getDiarySentenceListAll({ diaryId: diaryPreviewDiaryId.value })
  if (sentRes.code === 0) {
    const lang = String(diarySentencePreviewLang.value || 'zh')
    diarySentencePreviewTable.value = (Array.isArray(sentRes.data) ? sentRes.data : []).map((item) => {
      const translateObj = parseI18nObject(item?.translate)
      const normalizedTranslate = typeof item?.translate === 'string'
        ? item.translate
        : JSON.stringify(translateObj)
      return {
        ...item,
        translate: normalizedTranslate,
        translateObj,
        translateText: translateObj[lang] || translateObj.zh || translateObj.en || translateObj.mn || ''
      }
    })
    diarySentencePreviewOriginalMap.value = buildDiarySentencePreviewSnapshotMap(diarySentencePreviewTable.value)
    diaryPendingSentenceChanges.value = []
  }

  diarySkipAutoSelectAllWords.value = true
  await loadDiaryPreviewKeywords()
  diarySkipAutoSelectAllWords.value = false

  // 只保留本次实际更新高亮的单词的勾选状态
  if (diaryAllWordsTableRef.value) {
    diaryAllWordsTableRef.value.clearSelection()
    const idSet = new Set(allIds)
    const rowsToSelect = allDiarySubtitleWords.value.filter(k => idSet.has(k.wordId))
    rowsToSelect.forEach(row => {
      diaryAllWordsTableRef.value.toggleRowSelection(row, true)
    })
  }
  diaryAllWordsSelectedIds.value = [...allIds]
}

// 切换到"全部单词"tab 时自动选中已在词库的单词
const autoSelectDiaryAllMatchedWords = () => {
  if (diarySkipAutoSelectAllWords.value) return
  if (!diaryAllWordsTableRef.value) return
  setTimeout(() => {
    if (!diaryAllWordsTableRef.value) return
    diaryAllWordsTableRef.value.clearSelection()
    const rows = allDiarySubtitleWords.value.filter(k => k.isHighlighted)
    rows.forEach(row => {
      diaryAllWordsTableRef.value.toggleRowSelection(row, true)
    })
  }, 60)
}

watch(
  () => diarySentencePreviewTab.value,
  (newTab) => {
    if (newTab !== 'allWords') return
    autoSelectDiaryAllMatchedWords()
  }
)
watch(
  () => diaryPreviewKeywordList.value,
  () => {
    if (diarySentencePreviewTab.value !== 'allWords') return
    autoSelectDiaryAllMatchedWords()
  }
)

const computeDiaryPendingChanges = () => {
  const changes = []
  for (const item of diarySentencePreviewTable.value) {
    const rowId = Number(item?.id || item?.ID || 0)
    if (!rowId) continue
    const orig = diarySentencePreviewOriginalMap.value[rowId]
    if (!orig) continue
    const currentTranslateObj = JSON.stringify(parseI18nObject(item?.translateObj || item?.translate))
    const currentEnglish = String(item?.english || '').trim()
    if (orig.startTime !== Number(item?.startTime || 0) ||
        orig.endTime !== Number(item?.endTime || 0) ||
        orig.english !== currentEnglish ||
        orig.translateObj !== currentTranslateObj) {
      changes.push({
        id: rowId,
        startTime: Number(item?.startTime || 0),
        endTime: Number(item?.endTime || 0),
        english: currentEnglish,
        translate: typeof item.translate === 'string' ? item.translate : JSON.stringify(item.translate || {})
      })
    }
  }
  return changes
}

const diarySaveSentencePreview = async () => {
  const changes = computeDiaryPendingChanges()
  if (changes.length === 0) return
  diarySentenceSaving.value = true
  try {
    const res = await updateDiarySentenceList({
      diaryId: diaryPreviewDiaryId.value,
      sentences: changes
    })
    if (res.code === 0) {
      ElMessage.success(`保存成功，更新了 ${changes.length} 条`)
      diarySentencePreviewOriginalMap.value = buildDiarySentencePreviewSnapshotMap(diarySentencePreviewTable.value)
      diaryPendingSentenceChanges.value = []
    }
  } catch (e) {
    ElMessage.error('保存失败')
  }
  diarySentenceSaving.value = false
}

const diaryResetSentencePreviewChanges = () => {
  diarySentencePreviewTable.value = diarySentencePreviewTable.value.map((item) => {
    const rowId = Number(item?.id || item?.ID || 0)
    const orig = diarySentencePreviewOriginalMap.value[rowId]
    if (!orig) return item
    const translateObj = parseI18nObject(orig.translateObj)
    const lang = String(diarySentencePreviewLang.value || 'zh')
    return {
      ...item,
      startTime: orig.startTime,
      endTime: orig.endTime,
      english: orig.english,
      translate: orig.translateObj,
      translateObj,
      translateText: translateObj[lang] || translateObj.zh || translateObj.en || translateObj.mn || ''
    }
  })
  diaryPendingSentenceChanges.value = []
}

const handleDiarySentencePreviewBeforeClose = (done) => {
  const changes = computeDiaryPendingChanges()
  if (changes.length > 0) {
    ElMessageBox.confirm('有未保存的改动，确定关闭吗？', '提示', {
      confirmButtonText: '确定关闭',
      cancelButtonText: '继续编辑',
      type: 'warning'
    }).then(() => {
      diaryPendingSentenceChanges.value = []
      done()
    }).catch(() => {})
  } else {
    done()
  }
}

const diaryCloseSentencePreview = () => {
  diarySentencePreviewVisible.value = false
}

const previewDiarySubtitles = async (row) => {
  const diaryId = Number(row?.ID || 0)
  if (!diaryId) {
    ElMessage.warning('日记ID无效')
    return
  }

  const res = await getDiarySentenceListAll({ diaryId })
  if (res.code !== 0) {
    return
  }

  diaryPreviewDiaryId.value = diaryId
  diaryPreviewDiaryName.value = formatI18nText(row?.name)
  diarySentencePreviewTable.value = (Array.isArray(res.data) ? res.data : []).map((item) => {
    const translateObj = parseI18nObject(item?.translate)
    const lang = String(diarySentencePreviewLang.value || 'zh')
    // 规范化 translate 为字符串，避免对象/字符串类型不一致导致改动计数错误
    const normalizedTranslate = typeof item?.translate === 'string'
      ? item.translate
      : JSON.stringify(translateObj)
    return {
      ...item,
      translate: normalizedTranslate,
      translateObj,
      translateText: translateObj[lang] || translateObj.zh || translateObj.en || translateObj.mn || ''
    }
  })
  diarySentencePreviewOriginalMap.value = buildDiarySentencePreviewSnapshotMap(diarySentencePreviewTable.value)
  diarySentencePreviewTab.value = 'sentences'
  diarySentencePreviewVisible.value = true

  // 预加载关键词列表
  diaryPreviewKeywordList.value = []
  diaryPreviewSelectedKeywordIds.value = []
  loadDiaryPreviewKeywords()
}

// 字幕解析
const diarySubtitleSectionRef = ref(null)
const diarySubtitleForm = ref({ diaryId: undefined, englishSubtitleUrl: '', englishSubtitleUrlUk: '', translationSubtitleMap: {} })
const diarySubtitleFormKey = ref(0)
const diaryKeywordScanning = ref(false)
const diaryKeywordScanDone = ref(false)
const diaryKeywordScanResults = ref([])
const diarySelectedKeywordIds = ref([])
const diaryKeywordTableRef = ref(null)

const diarySubtitleUploadFolder = computed(() => {
  // 日记字幕上传目录
  return 'english-learn/diary/subtitle'
})

const subtitleLanguages = computed(() => {
  return managedLanguages.value.filter((lang) => String(lang.code || '').toLowerCase() !== 'en')
})

const openDiarySubtitleParser = async (row) => {
  // 通过 findDiary 获取完整日记数据（列表API可能不包含字幕URL字段）
  let diaryDetail = row
  try {
    const detailRes = await findDiary({ ID: row.ID })
    if (detailRes.code === 0 && detailRes.data) {
      diaryDetail = detailRes.data
    }
  } catch (e) {
    // 忽略，使用列表数据兜底
  }

  // 回填已有字幕文件路径
  let translationMap = {}
  try {
    if (diaryDetail.translationSubtitleUrls) {
      const parsed = typeof diaryDetail.translationSubtitleUrls === 'string'
        ? JSON.parse(diaryDetail.translationSubtitleUrls)
        : diaryDetail.translationSubtitleUrls
      if (parsed && typeof parsed === 'object') {
        translationMap = { ...parsed }
      }
    }
  } catch (e) {
    // JSON 解析失败忽略
  }
  diarySubtitleForm.value = {
    diaryId: row.ID,
    englishSubtitleUrl: diaryDetail.englishSubtitleUrl || '',
    englishSubtitleUrlUk: diaryDetail.englishSubtitleUrlUk || '',
    translationSubtitleMap: translationMap
  }
  // 强制触发 FileUploadWithDir 组件更新（key 变化触发重新渲染）
  diarySubtitleFormKey.value++
  diaryKeywordScanDone.value = false
  diaryKeywordScanResults.value = []
  diarySelectedKeywordIds.value = []
  activeTab.value = 'diary'
  setTimeout(() => {
    diarySubtitleSectionRef.value?.scrollIntoView({ behavior: 'smooth' })
  }, 200)

  // 如果已有字幕，加载完整词库并预选已高亮的单词
  try {
    const kwRes = await getDiaryKeywords({ diaryId: row.ID })
    if (kwRes.code === 0 && kwRes.data && kwRes.data.length > 0) {
      diaryKeywordScanResults.value = kwRes.data.map(item => ({
        word: item.word,
        wordId: item.wordId,
        matched: !!item.matched,
        count: item.count || 0
      }))
      // 预选已高亮的单词（matched=true 的）
      const highlightedIds = diaryKeywordScanResults.value
        .filter(k => k.matched)
        .map(k => k.wordId)
      diarySelectedKeywordIds.value = highlightedIds
      diaryKeywordScanDone.value = true
      // 自动选中 matched 的表格行
      nextTick(() => {
        if (diaryKeywordTableRef.value) {
          const matchedRows = diaryKeywordScanResults.value.filter(k => k.matched)
          matchedRows.forEach(rowItem => {
            diaryKeywordTableRef.value.toggleRowSelection(rowItem, true)
          })
        }
      })
      if (highlightedIds.length > 0) {
        ElMessage.info(`该日记已有字幕，字幕共 ${diaryKeywordScanResults.value.length} 个单词，其中 ${highlightedIds.length} 个匹配词库`)
      }
    }
  } catch (e) {
    // 忽略加载错误
  }
}

const handleDiaryScanKeywords = async () => {
  if (!diarySubtitleForm.value.diaryId) {
    ElMessage.warning('请选择目标日记')
    return
  }
  if (!diarySubtitleForm.value.englishSubtitleUrl) {
    ElMessage.warning('请先上传英文字幕文件')
    return
  }
  const previousKeywordIds = [...diarySelectedKeywordIds.value]
  diaryKeywordScanning.value = true
  diaryKeywordScanDone.value = false
  diaryKeywordScanResults.value = []
  diarySelectedKeywordIds.value = []
  const res = await scanDiaryKeywords({
    diaryId: diarySubtitleForm.value.diaryId,
    englishSubtitleUrl: diarySubtitleForm.value.englishSubtitleUrl
  })
  diaryKeywordScanning.value = false
  if (res.code === 0) {
    // 后端返回字幕中出现的所有单词，matched=true 表示在词库中匹配到
    diaryKeywordScanResults.value = res.data || []
    diaryKeywordScanDone.value = true
    // 合并：扫描匹配的 + 之前已选的（且仍在结果中的）
    const matchedIds = diaryKeywordScanResults.value.filter(k => k.matched).map(k => k.wordId)
    const previousInResults = previousKeywordIds.filter(id =>
      diaryKeywordScanResults.value.some(k => k.wordId === id)
    )
    const mergedIds = [...new Set([...matchedIds, ...previousInResults])]
    diarySelectedKeywordIds.value = mergedIds
    nextTick(() => {
      if (diaryKeywordTableRef.value) {
        const rowsToSelect = diaryKeywordScanResults.value.filter(k => mergedIds.includes(k.wordId))
        rowsToSelect.forEach(row => diaryKeywordTableRef.value.toggleRowSelection(row, true))
      }
    })
    ElMessage.success(`扫描完成，字幕共 ${diaryKeywordScanResults.value.length} 个单词，其中 ${matchedIds.length} 个匹配词库`)
  }
}

const handleDiaryKeywordSelectionChange = (rows) => {
  diarySelectedKeywordIds.value = rows.map(r => r.wordId).filter(Boolean)
}

const diarySelectAllKeywords = (select) => {
  if (!diaryKeywordTableRef.value) return
  if (select) {
    const matchedRows = diaryKeywordScanResults.value.filter(k => k.matched)
    matchedRows.forEach(row => diaryKeywordTableRef.value.toggleRowSelection(row, true))
  } else {
    diaryKeywordTableRef.value.clearSelection()
  }
}

const handleDiaryConfirmParse = async () => {
  if (!diarySubtitleForm.value.diaryId) {
    ElMessage.warning('请选择目标日记')
    return
  }
  const translationSubtitle = []
  for (const lang of subtitleLanguages.value) {
    const url = diarySubtitleForm.value.translationSubtitleMap[lang.code]
    if (url) {
      translationSubtitle.push({ language: lang.code, subtitleUrl: url })
    }
  }
  const res = await parseDiarySubtitleFiles({
    diaryId: diarySubtitleForm.value.diaryId,
    englishSubtitleUrl: diarySubtitleForm.value.englishSubtitleUrl,
    englishSubtitleUrlUk: diarySubtitleForm.value.englishSubtitleUrlUk,
    translationSubtitle,
    keywordIds: diarySelectedKeywordIds.value
  })
  if (res.code === 0) {
    ElMessage.success('字幕解析入库成功')
    diarySubtitleForm.value = { diaryId: undefined, englishSubtitleUrl: '', englishSubtitleUrlUk: '', translationSubtitleMap: {} }
    diarySubtitleFormKey.value++
    diaryKeywordScanDone.value = false
    diaryKeywordScanResults.value = []
    diarySelectedKeywordIds.value = []
    // 刷新日记列表，确保表格中显示最新字幕URL
    loadDiaryList()
    loadDiaryOptions()
  }
}

// ===== 授权 =====
const entitlementQuery = ref({ page: 1, pageSize: 10, userId: undefined, resourceType: undefined })
const entitlementTable = ref([])
const entitlementTotal = ref(0)
const entitlementDialogVisible = ref(false)
const entitlementForm = ref({
  userId: undefined,
  resourceType: 'diary',
  resourceId: undefined,
  expireAt: '',
  remark: ''
})

const resourceTypeLabel = (type) => {
  return { english_category: '分类', video_series: '剧集', video_episode: '单集', diary: '日记' }[type] || type
}

const resourceOptions = computed(() => {
  if (entitlementForm.value.resourceType === 'english_category') {
    return diaryCategoryOptions.value.map((item) => ({ value: item.ID, label: `${item.ID} - ${formatI18nText(item.name)}` }))
  }
  if (entitlementForm.value.resourceType === 'diary') {
    return diaryOptions.value.map((item) => ({ value: item.ID, label: `${item.ID} - ${formatI18nText(item.name)}` }))
  }
  return []
})

const loadEntitlementList = async () => {
  const res = await getEntitlementList(entitlementQuery.value)
  if (res.code === 0 && res.data) {
    entitlementTable.value = res.data.list || []
    entitlementTotal.value = res.data.total || 0
  }
}

const handleEntitlementSearch = () => { entitlementQuery.value.page = 1; loadEntitlementList() }
const resetEntitlementSearch = () => {
  entitlementQuery.value = { page: 1, pageSize: 10, userId: undefined, resourceType: undefined }
  loadEntitlementList()
}
const handleEntitlementPageChange = (page) => { entitlementQuery.value.page = page; loadEntitlementList() }
const handleEntitlementSizeChange = (size) => { entitlementQuery.value.pageSize = size; entitlementQuery.value.page = 1; loadEntitlementList() }

const openEntitlementDialog = () => {
  entitlementForm.value = { userId: undefined, resourceType: 'diary', resourceId: undefined, expireAt: '', remark: '' }
  entitlementDialogVisible.value = true
}

const submitEntitlement = async () => {
  const res = await grantEntitlement(entitlementForm.value)
  if (res.code === 0) {
    ElMessage.success('授权成功')
    entitlementDialogVisible.value = false
    loadEntitlementList()
  }
}

const removeEntitlement = (row) => {
  ElMessageBox.confirm('确定要撤销该授权吗？', '提示', { type: 'warning' }).then(async () => {
    const res = await revokeEntitlement({ userId: row.userId, resourceType: row.resourceType, resourceId: row.resourceId })
    if (res.code === 0) {
      ElMessage.success('撤销成功')
      loadEntitlementList()
    }
  })
}

const loadAll = () => {
  loadDiaryCategoryList()
  loadDiaryCategoryOptions()
  loadDiaryList()
  loadDiaryOptions()
  loadDiaryTagOptions()
  loadEntitlementList()
}

onMounted(async () => {
  try {
    const res = await getLanguageList()
    if (res.code === 0 && res.data) {
      managedLanguages.value = Array.isArray(res.data) ? res.data : (res.data.list || [])
    }
  } catch (e) { /* ignore */ }
  loadAll()
})
</script>

<style scoped>
.lang-switch-row {
  margin-top: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 12px;
}

.toolbar-tip {
  color: #909399;
  font-size: 13px;
}

.subtitle-box {
  margin-top: 32px;
}

.subtitle-preview-toolbar {
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.mt-2 { margin-top: 8px; }
.w-full { width: 100%; }

.keyword-scan-result {
  margin-top: 16px;
}

.keyword-toolbar {
  margin-bottom: 12px;
  display: flex;
  gap: 8px;
}
</style>