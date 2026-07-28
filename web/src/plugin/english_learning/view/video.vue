<template>
  <div>
    <div class="gva-search-box">
      <el-alert
        title="英语学习运营台（视频端）：支持视频分类、剧集、单集、字幕文件解析和用户资源授权管理。"
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
        <el-tab-pane label="视频分类" name="videoCategory">
          <div class="toolbar-row">
            <span class="toolbar-tip">用于维护视频大类（如影视英语、情景口语等）与存储目录键</span>
            <el-button type="primary" icon="plus" @click="openVideoCategoryDialog()">新增视频分类</el-button>
          </div>

          <el-table :data="videoCategoryTable" row-key="ID" border>
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
                <el-button type="primary" link icon="edit" @click="openVideoCategoryDialog(scope.row)">编辑</el-button>
                <el-button type="danger" link icon="delete" @click="removeVideoCategory(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="videoCategoryQuery.page"
              :page-size="videoCategoryQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="videoCategoryTotal"
              @current-change="handleVideoCategoryPageChange"
              @size-change="handleVideoCategorySizeChange"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane label="剧集管理" name="series">
          <div class="toolbar-row">
            <el-form :inline="true" :model="seriesQuery">
              <el-form-item label="所属分类">
                <el-select v-model="seriesQuery.categoryId" clearable placeholder="全部" style="width: 240px">
                  <el-option
                    v-for="item in videoCategoryOptions"
                    :key="item.ID"
                    :label="formatI18nText(item.name)"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="剧集关键词">
                <el-input v-model="seriesQuery.keyword" clearable placeholder="名称模糊查询" style="width: 220px" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="handleSeriesSearch">查询</el-button>
                <el-button icon="refresh" @click="resetSeriesSearch">重置</el-button>
              </el-form-item>
            </el-form>
            <el-button type="primary" icon="plus" @click="openSeriesDialog()">新增剧集</el-button>
          </div>

          <el-table :data="seriesTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column label="所属分类" min-width="180">
              <template #default="scope">
                {{ videoCategoryNameMap[scope.row.categoryId] || `#${scope.row.categoryId}` }}
              </template>
            </el-table-column>
            <el-table-column label="剧集名称" min-width="220">
              <template #default="scope">
                {{ formatI18nText(scope.row.name) }}
              </template>
            </el-table-column>
            <el-table-column prop="coverUrl" label="封面地址" min-width="240" show-overflow-tooltip />
            <el-table-column prop="price" label="价格" width="120" />
            <el-table-column label="首页展示" width="120">
              <template #default="scope">
                <el-tag :type="scope.row.showHome === false ? 'info' : 'success'">
                  {{ scope.row.showHome === false ? '隐藏' : '展示' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="会员限制" width="120">
              <template #default="scope">
                <el-tag :type="scope.row.needVip ? 'warning' : 'success'">{{ scope.row.needVip ? '需要会员' : '普通可见' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="viewCount" label="浏览量" width="100" />
            <el-table-column prop="userCount" label="观看人数" width="110" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openSeriesDialog(scope.row)">编辑</el-button>
                <el-button type="danger" link icon="delete" @click="removeSeries(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="seriesQuery.page"
              :page-size="seriesQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="seriesTotal"
              @current-change="handleSeriesPageChange"
              @size-change="handleSeriesSizeChange"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane label="单集与字幕" name="episode">
          <div class="toolbar-row">
            <el-form :inline="true" :model="episodeQuery">
              <el-form-item label="所属剧集">
                <el-select v-model="episodeQuery.seriesId" clearable filterable placeholder="全部" style="width: 280px">
                  <el-option
                    v-for="item in seriesOptions"
                    :key="item.ID"
                    :label="formatI18nText(item.name)"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="单集关键词">
                <el-input v-model="episodeQuery.keyword" clearable placeholder="名称模糊查询" style="width: 220px" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="handleEpisodeSearch">查询</el-button>
                <el-button icon="refresh" @click="resetEpisodeSearch">重置</el-button>
              </el-form-item>
            </el-form>
            <el-button type="primary" icon="plus" @click="openEpisodeDialog()">新增单集</el-button>
          </div>

          <el-table :data="episodeTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column label="所属剧集" min-width="220">
              <template #default="scope">
                {{ seriesNameMap[scope.row.seriesId] || `#${scope.row.seriesId}` }}
              </template>
            </el-table-column>
            <el-table-column label="单集名称" min-width="220">
              <template #default="scope">
                {{ formatI18nText(scope.row.name) }}
              </template>
            </el-table-column>
            <el-table-column prop="videoUrl" label="视频地址" min-width="260" show-overflow-tooltip />
            <el-table-column prop="trialPercent" label="试看比例(%)" width="120" />
            <el-table-column prop="sort" label="排序" width="100" />
            <el-table-column label="操作" width="320" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openEpisodeDialog(scope.row)">编辑</el-button>
                <el-button type="success" link @click="openSubtitleParser(scope.row)">字幕解析</el-button>
                <el-button type="warning" link @click="previewEpisodeSubtitles(scope.row)">查看字幕</el-button>
                <el-button type="danger" link icon="delete" @click="removeEpisode(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="episodeQuery.page"
              :page-size="episodeQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="episodeTotal"
              @current-change="handleEpisodePageChange"
              @size-change="handleEpisodeSizeChange"
            />
          </div>

          <div ref="subtitleSectionRef" class="subtitle-box">
            <el-divider content-position="left">字幕文件解析入库</el-divider>
            <el-form label-width="140px">
              <el-form-item label="目标单集">
                <el-select v-model="subtitleForm.episodeId" filterable placeholder="请选择单集" style="width: 480px">
                  <el-option
                    v-for="item in episodeOptions"
                    :key="item.ID"
                    :label="`${seriesNameMap[item.seriesId] || '-'} / ${formatI18nText(item.name)}`"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="英文字幕文件(必填)">
                <FileUploadWithDir v-model="subtitleForm.englishSubtitleUrl" :default-folder="subtitleUploadFolder" accept=".srt,.vtt,.ass" />
              </el-form-item>

              <el-form-item
                v-for="lang in subtitleLanguages"
                :key="lang.code"
                :label="`${lang.name || lang.code} 字幕文件`"
              >
                <FileUploadWithDir v-model="subtitleForm.translationSubtitleMap[lang.code]" :default-folder="subtitleUploadFolder" accept=".srt,.vtt,.ass" />
              </el-form-item>

              <el-form-item>
                <el-button type="warning" :loading="keywordScanning" @click="handleScanKeywords">
                  {{ keywordScanning ? '正在扫描...' : '步骤1: 扫描关键词' }}
                </el-button>
                <el-button
                  v-if="keywordScanDone"
                  type="primary"
                  style="margin-left: 12px"
                  :disabled="selectedKeywordIds.length === 0"
                  @click="handleConfirmParse"
                >
                  步骤2: 确认解析入库（已选 {{ selectedKeywordIds.length }} 个词）
                </el-button>
              </el-form-item>
            </el-form>

            <!-- 关键词扫描结果 -->
            <div v-if="keywordScanDone" class="keyword-scan-result">
              <el-divider content-position="left">
                关键词扫描结果（共 {{ keywordScanResults.length }} 个单词，已选 {{ selectedKeywordIds.length }} 个）
              </el-divider>
              <div class="keyword-toolbar">
                <el-button size="small" @click="selectAllKeywords(true)">全选已匹配</el-button>
                <el-button size="small" @click="selectAllKeywords(false)">取消全选</el-button>
              </div>
              <el-table
                ref="keywordTableRef"
                :data="keywordScanResults"
                row-key="word"
                border
                max-height="400"
                @selection-change="handleKeywordSelectionChange"
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

        <el-tab-pane label="用户授权" name="entitlement">
          <div class="toolbar-row">
            <el-form :inline="true" :model="entitlementQuery">
              <el-form-item label="用户ID">
                <el-input-number v-model="entitlementQuery.userId" :min="1" controls-position="right" />
              </el-form-item>
              <el-form-item label="资源类型">
                <el-select v-model="entitlementQuery.resourceType" clearable placeholder="全部" style="width: 220px">
                  <el-option label="分类" value="english_category" />
                  <el-option label="剧集" value="video_series" />
                  <el-option label="单集" value="video_episode" />
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

    <el-dialog v-model="videoCategoryDialogVisible" :title="videoCategoryDialogMode === 'create' ? '新增视频分类' : '编辑视频分类'" width="760px" :close-on-click-modal="false" :close-on-press-escape="false">
      <el-form :model="videoCategoryForm" label-width="120px">
        <el-form-item label="分类名称">
          <MultiLangEditor
            :model="videoCategoryForm.nameI18n"
            title="视频分类多语言名称"
            input-type="input"
            :rows="2"
            :use-tabs="true"
          />
        </el-form-item>
        <el-form-item label="存储目录键">
          <el-input v-model="videoCategoryForm.storageKey" placeholder="例如: liblib（用于视频/字幕目录 english-learn/video/liblib）" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="videoCategoryForm.sort" :min="0" />
        </el-form-item>
        <el-form-item label="首页展示">
          <el-switch v-model="videoCategoryForm.showHome" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="videoCategoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitVideoCategory">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="seriesDialogVisible" :title="seriesDialogMode === 'create' ? '新增剧集' : '编辑剧集'" width="860px" :close-on-click-modal="false" :close-on-press-escape="false">
      <el-form :model="seriesForm" label-width="120px">
        <el-form-item label="所属分类">
          <el-select v-model="seriesForm.categoryId" placeholder="请选择分类" style="width: 100%">
            <el-option
              v-for="item in videoCategoryOptions"
              :key="item.ID"
              :label="formatI18nText(item.name)"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="剧集名称">
          <MultiLangEditor
            :model="seriesForm.nameI18n"
            title="剧集名称多语言"
            input-type="input"
            :rows="2"
            :use-tabs="true"
          />
        </el-form-item>
        <el-form-item label="封面地址">
          <FileUploadWithDir v-model="seriesForm.coverUrl" default-folder="english-learn/video/cover" accept="image/*" />
        </el-form-item>
        <el-form-item label="价格">
          <el-input-number v-model="seriesForm.price" :min="0" :precision="2" :step="1" />
        </el-form-item>
        <el-form-item label="需要会员">
          <el-switch v-model="seriesForm.needVip" />
        </el-form-item>
        <el-form-item label="首页展示">
          <el-switch v-model="seriesForm.showHome" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="seriesDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitSeries">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="episodeDialogVisible" :title="episodeDialogMode === 'create' ? '新增单集' : '编辑单集'" width="860px" :close-on-click-modal="false" :close-on-press-escape="false">
      <el-form :model="episodeForm" label-width="120px">
        <el-form-item label="所属剧集">
          <el-select v-model="episodeForm.seriesId" filterable placeholder="请选择剧集" style="width: 100%">
            <el-option
              v-for="item in seriesOptions"
              :key="item.ID"
              :label="formatI18nText(item.name)"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="单集名称">
          <MultiLangEditor
            :model="episodeForm.nameI18n"
            title="单集名称多语言"
            input-type="input"
            :rows="2"
            :use-tabs="true"
          />
        </el-form-item>
        <el-form-item label="切片模式">
          <el-switch
            v-model="hlsMode"
            active-text="HLS切片"
            inactive-text="普通上传"
            @change="onHlsModeChange"
          />
          <span v-if="hlsMode && ffmpegChecking" style="margin-left:12px;color:#409eff;">
            ⏳ 正在检测 FFmpeg...
          </span>
          <el-tag v-else-if="hlsMode && ffmpegReady" type="success" style="margin-left:12px;">FFmpeg 就绪</el-tag>
          <el-tag v-else-if="hlsMode && !ffmpegReady" type="danger" style="margin-left:12px;">FFmpeg 不可用</el-tag>
        </el-form-item>
        <el-form-item v-if="!hlsMode" label="视频地址">
          <FileUploadWithDir v-model="episodeForm.videoUrl" :default-folder="episodeUploadFolder" accept="video/*" />
        </el-form-item>
        <el-form-item v-else label="HLS切片">
          <div class="hls-folder-bar">
            <span class="hls-folder-label">切片目录：</span>
            <el-input
              v-model="hlsFolder"
              size="small"
              class="hls-folder-input"
              placeholder="如 didi，留空自动生成"
              clearable
            />
          </div>
          <div class="hls-folder-bar" style="margin-top: 8px;">
            <span class="hls-folder-label">视频文件：</span>
            <input
              ref="hlsFileInput"
              type="file"
              accept="video/*"
              style="flex:1;"
              @change="onHlsFileChange"
            />
            <el-button
              type="primary"
              style="margin-left: 8px;"
              :loading="hlsSlicing"
              :disabled="!ffmpegReady || !hlsFile"
              @click="onHlsSliceClick"
            >
              {{ hlsSlicing ? '正在切片上传...' : '开始切片' }}
            </el-button>
          </div>
          <div v-if="hlsProgress" style="margin-top:8px;color:#409eff;">{{ hlsProgress }}</div>
        </el-form-item>
        <el-form-item label="试看比例(%)">
          <el-input-number v-model="episodeForm.trialPercent" :min="1" :max="100" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="episodeForm.sort" :min="0" />
        </el-form-item>
        <el-form-item label="视频标签">
          <el-select
            v-model="episodeForm.tagIds"
            multiple
            filterable
            clearable
            placeholder="请选择视频标签"
            style="width: 100%"
          >
            <el-option
              v-for="item in videoTagOptions"
              :key="item.ID"
              :label="item.name"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="episodeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEpisode" :loading="hlsSlicing">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="sentencePreviewVisible"
      :title="`字幕预览 - ${sentencePreviewEpisodeName || ''}`"
      width="1100px"
      :before-close="handleSentencePreviewBeforeClose"
    >
      <el-tabs v-model="sentencePreviewTab">
        <el-tab-pane label="字幕句子" name="sentences">
          <div class="subtitle-preview-toolbar">
            <span>字幕语言：</span>
            <el-select v-model="sentencePreviewLang" style="width: 140px">
              <el-option
                v-for="lang in displayLanguageOptions"
                :key="lang.value"
                :label="lang.label"
                :value="lang.value"
              />
            </el-select>
            <el-tag type="warning" effect="light">已改动 {{ pendingSentenceChanges.length }} 条</el-tag>
          </div>
          <el-table :data="sentencePreviewTable" border max-height="500">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column label="开始(s)" width="120">
              <template #default="scope">
                <el-input-number v-model="scope.row.startTime" :min="0" :step="0.1" controls-position="right" />
              </template>
            </el-table-column>
            <el-table-column label="结束(s)" width="120">
              <template #default="scope">
                <el-input-number v-model="scope.row.endTime" :min="0" :step="0.1" controls-position="right" />
              </template>
            </el-table-column>
            <el-table-column label="英文句子" min-width="320">
              <template #default="scope">
                <el-input v-model="scope.row.english" type="textarea" :rows="2" />
              </template>
            </el-table-column>
            <el-table-column label="翻译" min-width="320">
              <template #default="scope">
                <el-input v-model="scope.row.translateText" type="textarea" :rows="2" @input="syncSentenceTranslate(scope.row)" />
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="重点单词" name="keywords">
          <div class="subtitle-preview-toolbar">
            <span>当前高亮的重点单词（共 {{ episodeKeywordList.length }} 个）</span>
            <el-button size="small" type="primary" @click="loadEpisodeKeywords">刷新列表</el-button>
          </div>
          <el-table
            ref="previewKeywordTableRef"
            :data="episodeKeywordList"
            row-key="word"
            border
            max-height="420"
            @selection-change="handlePreviewKeywordSelectionChange"
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
            <el-button size="small" @click="selectAllPreviewKeywords(true)">全选已匹配</el-button>
            <el-button size="small" @click="selectAllPreviewKeywords(false)">取消全选</el-button>
            <el-button
              type="warning"
              size="small"
              style="margin-left: 16px"
              :loading="previewRehighlighting"
              :disabled="previewSelectedKeywordIds.length === 0"
              @click="handleRehighlightPreview"
            >
              更新高亮（已选 {{ previewSelectedKeywordIds.length }} 个词）
            </el-button>
          </div>
        </el-tab-pane>

        <el-tab-pane label="全部单词" name="allWords">
          <div class="subtitle-preview-toolbar">
            <span>当前视频字幕中的所有单词（共 {{ allSubtitleWords.length }} 个），勾选后点击"更新高亮"即可添加为重点单词</span>
          </div>
          <el-table
            ref="allWordsTableRef"
            :data="allSubtitleWords"
            row-key="word"
            border
            max-height="420"
            @selection-change="handleAllWordsSelectionChange"
          >
            <el-table-column type="selection" width="50" />
            <el-table-column prop="word" label="单词" width="180" />
            <el-table-column label="已在词库" width="100">
              <template #default="scope">
                <el-tag v-if="scope.row.matched" type="success" size="small">是(ID:{{ scope.row.wordId }})</el-tag>
                <el-tag v-else type="info" size="small">否</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="count" label="出现次数" width="100" sortable />
          </el-table>
          <div style="margin-top: 12px">
            <el-button size="small" @click="selectAllAllWords(true)">全选已匹配</el-button>
            <el-button size="small" @click="selectAllAllWords(false)">取消全选</el-button>
            <el-button
              type="warning"
              size="small"
              style="margin-left: 16px"
              :loading="allWordsRehighlighting"
              :disabled="allWordsSelectedIds.length === 0"
              @click="handleRehighlightAllWords"
            >
              更新高亮（已选 {{ allWordsSelectedIds.length }} 个词）
            </el-button>
          </div>
        </el-tab-pane>
      </el-tabs>

      <template #footer>
        <el-button :disabled="pendingSentenceChanges.length === 0" @click="resetSentencePreviewChanges">重置改动</el-button>
        <el-button type="primary" :loading="sentenceSaving" :disabled="pendingSentenceChanges.length === 0" @click="saveSentencePreview">保存修改</el-button>
        <el-button @click="closeSentencePreview">关闭</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="entitlementDialogVisible" title="新增资源授权" width="760px" :close-on-click-modal="false" :close-on-press-escape="false">
      <el-form :model="entitlementForm" label-width="120px">
        <el-form-item label="用户ID">
          <el-input-number v-model="entitlementForm.userId" :min="1" controls-position="right" />
        </el-form-item>
        <el-form-item label="资源类型">
          <el-select v-model="entitlementForm.resourceType" style="width: 100%">
            <el-option label="分类" value="english_category" />
            <el-option label="剧集" value="video_series" />
            <el-option label="单集" value="video_episode" />
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
  </div>
</template>

<script setup>
  import { computed, onMounted, ref, watch } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'
  import FileUploadWithDir from '@/components/FileUploadWithDir/index.vue'
  import { getLanguageList } from '@/api/client/language'
  import {
    createVideoCategory,
    createVideoEpisode,
    createVideoSeries,
    findVideoEpisode,
    deleteVideoCategory,
    deleteVideoEpisode,
    deleteVideoSeries,
    getEntitlementList,
    getEpisodeKeywords,
    getEpisodeSubtitles,
    batchCreateWords,
    getVideoCategoryList,
    getVideoEpisodeList,
    getVideoSentenceList,
    getVideoSeriesList,
    grantEntitlement,
    parseSubtitleFiles,
    rehighlightSentences,
    revokeEntitlement,
    scanKeywords,
    updateVideoSentenceList,
    updateVideoCategory,
    updateVideoEpisode,
    updateVideoSeries,
    checkFfmpeg,
    sliceVideoEpisode
  } from '../api/english'
  import { getVideoTagList } from '@/api/client/videoTag'

  defineOptions({
    name: 'EnglishLearningVideo'
  })

  const activeTab = ref('videoCategory')
  const displayLang = ref('zh')

  const normalizeI18nObject = (raw) => {
    if (!raw) return { zh: '' }
    if (typeof raw === 'object') {
      return Object.keys(raw).length > 0 ? { ...raw } : { zh: '' }
    }
    const text = String(raw).trim()
    if (!text) return { zh: '' }
    try {
      const parsed = JSON.parse(text)
      if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
        return Object.keys(parsed).length > 0 ? { ...parsed } : { zh: '' }
      }
      return { zh: text }
    } catch (e) {
      return { zh: text }
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

  const formatI18nText = (raw) => {
    if (!raw) return ''
    try {
      const parsed = typeof raw === 'string' ? JSON.parse(raw) : raw
      if (typeof parsed === 'object' && parsed !== null) {
        const lang = String(displayLang.value || 'zh')
        return parsed[lang] || parsed.zh || parsed.en || parsed.mn || Object.values(parsed)[0] || ''
      }
      return String(raw)
    } catch (e) {
      return String(raw)
    }
  }

  const parseI18nObject = (raw) => {
    if (!raw) return {}
    if (typeof raw === 'object') {
      return { ...raw }
    }
    try {
      const parsed = JSON.parse(String(raw))
      return parsed && typeof parsed === 'object' && !Array.isArray(parsed) ? { ...parsed } : {}
    } catch (e) {
      return {}
    }
  }

  const normalizeStorageKey = (raw) => {
    return String(raw || '')
      .toLowerCase()
      .trim()
      .replace(/[^a-z0-9_-]+/g, '-')
      .replace(/^-+|-+$/g, '')
  }

  const videoCategoryQuery = ref({ page: 1, pageSize: 10 })
  const videoCategoryTable = ref([])
  const videoCategoryTotal = ref(0)

  const seriesQuery = ref({ page: 1, pageSize: 10, categoryId: undefined, keyword: '' })
  const seriesTable = ref([])
  const seriesTotal = ref(0)

  const episodeQuery = ref({ page: 1, pageSize: 10, seriesId: undefined, keyword: '' })
  const episodeTable = ref([])
  const episodeTotal = ref(0)

  const entitlementQuery = ref({ page: 1, pageSize: 10, userId: undefined, resourceType: undefined })
  const entitlementTable = ref([])
  const entitlementTotal = ref(0)

  const videoCategoryOptions = ref([])
  const seriesOptions = ref([])
  const episodeOptions = ref([])

  const managedLanguages = ref([])

  const displayLanguageOptions = computed(() => {
    if (managedLanguages.value.length === 0) {
      return [
        { value: 'zh', label: '中文' },
        { value: 'en', label: 'English' },
        { value: 'mn', label: 'Монгол' }
      ]
    }
    return managedLanguages.value.map((lang) => ({
      value: lang.code,
      label: lang.nativeName || lang.name || lang.code
    }))
  })

  const subtitleLanguages = computed(() => {
    return managedLanguages.value.filter((lang) => String(lang.code || '').toLowerCase() !== 'en')
  })

  const videoCategoryNameMap = computed(() => {
    const map = {}
    for (const item of videoCategoryOptions.value) {
      map[item.ID] = formatI18nText(item.name)
    }
    return map
  })

  const seriesNameMap = computed(() => {
    const map = {}
    for (const item of seriesOptions.value) {
      map[item.ID] = formatI18nText(item.name)
    }
    return map
  })

  const videoCategoryDialogVisible = ref(false)
  const videoCategoryDialogMode = ref('create')
  const videoCategoryForm = ref({
    ID: 0,
    nameI18n: { zh: '' },
    storageKey: '',
    showHome: true,
    sort: 0
  })

  const seriesDialogVisible = ref(false)
  const seriesDialogMode = ref('create')
  const seriesForm = ref({
    ID: 0,
    categoryId: undefined,
    nameI18n: { zh: '' },
    coverId: 0,
    coverUrl: '',
    price: 0,
    needVip: false,
    showHome: true
  })

  const episodeDialogVisible = ref(false)
  const episodeDialogMode = ref('create')
  const episodeForm = ref({
    ID: 0,
    seriesId: undefined,
    nameI18n: { zh: '' },
    videoUrl: '',
    trialPercent: 8,
    sort: 0,
    tagIds: []
  })

  const videoTagOptions = ref([])

  // HLS 切片相关
  const hlsMode = ref(false)
  const ffmpegReady = ref(false)
  const ffmpegChecking = ref(false)
  const hlsSlicing = ref(false)
  const hlsProgress = ref('')
  const hlsFolder = ref('')
  const hlsFile = ref(null)
  const hlsFileInput = ref(null)

  const subtitleForm = ref({
    episodeId: undefined,
    englishSubtitleUrl: '',
    translationSubtitleMap: {}
  })
  const subtitleSectionRef = ref(null)

  // 关键词扫描相关
  const keywordScanning = ref(false)
  const keywordScanDone = ref(false)
  const keywordScanResults = ref([])
  const selectedKeywordIds = ref([])
  const keywordTableRef = ref(null)

  const sentencePreviewVisible = ref(false)
  const sentencePreviewEpisodeName = ref('')
  const sentencePreviewTable = ref([])
  const sentencePreviewEpisodeId = ref(0)
  const sentencePreviewLang = ref('zh')
  const sentencePreviewTab = ref('sentences')
  const sentenceSaving = ref(false)
  const sentencePreviewOriginalMap = ref({})

  // 预览弹窗-关键词管理
  const episodeKeywordList = ref([])
  const previewSelectedKeywordIds = ref([])
  const previewRehighlighting = ref(false)
  const previewKeywordTableRef = ref(null)

  // 全部单词 tab - 选择管理
  const allWordsTableRef = ref(null)
  const allWordsSelectedIds = ref([])
  const allWordsRehighlighting = ref(false)

  const normalizeSentenceSnapshot = (item) => {
    const normalizeTranslateObj = (raw) => {
      const source = parseI18nObject(raw)
      const sorted = {}
      for (const key of Object.keys(source).sort()) {
        sorted[key] = String(source[key] ?? '').trim()
      }
      return sorted
    }
    return {
      english: String(item?.english || '').trim(),
      startTime: Number(item?.startTime || 0),
      endTime: Number(item?.endTime || 0),
      translateObj: normalizeTranslateObj(item?.translateObj || item?.translate)
    }
  }

  const buildSentencePreviewSnapshotMap = (rows) => {
    const snapshot = {}
    for (const row of rows || []) {
      const rowId = Number(row?.id || row?.ID || 0)
      if (!rowId) continue
      snapshot[rowId] = normalizeSentenceSnapshot(row)
    }
    return snapshot
  }

  const isSentenceRowChanged = (row) => {
    const rowId = Number(row?.id || row?.ID || 0)
    if (!rowId) return false
    const current = normalizeSentenceSnapshot(row)
    const original = sentencePreviewOriginalMap.value[rowId]
    if (!original) return true
    return (
      current.english !== original.english ||
      current.startTime !== original.startTime ||
      current.endTime !== original.endTime ||
      JSON.stringify(current.translateObj) !== JSON.stringify(original.translateObj)
    )
  }

  const pendingSentenceChanges = computed(() => {
    return sentencePreviewTable.value.filter((row) => isSentenceRowChanged(row))
  })

  // 全部单词：从当前字幕句子中提取所有单词并统计词频（用于"查看字幕-全部单词"tab）
  const allSubtitleWords = computed(() => {
    const wordCounter = new Map()
    const wordRegex = /[a-zA-Z]+/g
    // 收集已高亮的关键词信息用于匹配
    const highlightedWordMap = new Map()
    for (const kw of episodeKeywordList.value) {
      if (kw.matched) {
        highlightedWordMap.set(kw.word.toLowerCase(), kw.wordId)
      }
    }

    for (const row of sentencePreviewTable.value) {
      const text = String(row?.english || '').replace(/<[^>]+>/g, '') // 去除HTML标签
      const matches = text.matchAll(wordRegex)
      for (const m of matches) {
        const word = m[0].toLowerCase()
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
      wordId: highlightedWordMap.get(word) || 0
    }))
  })

  const entitlementDialogVisible = ref(false)
  const entitlementForm = ref({
    userId: undefined,
    resourceType: 'video_series',
    resourceId: undefined,
    expireAt: '',
    remark: ''
  })

  const resourceOptions = computed(() => {
    if (entitlementForm.value.resourceType === 'english_category') {
      return videoCategoryOptions.value.map((item) => ({
        value: item.ID,
        label: `${item.ID} - ${formatI18nText(item.name)}`
      }))
    }
    if (entitlementForm.value.resourceType === 'video_series') {
      return seriesOptions.value.map((item) => ({
        value: item.ID,
        label: `${item.ID} - ${formatI18nText(item.name)}`
      }))
    }
    return episodeOptions.value.map((item) => ({
      value: item.ID,
      label: `${item.ID} - ${seriesNameMap.value[item.seriesId] || '-'} / ${formatI18nText(item.name)}`
    }))
  })

  const getCategoryByID = (categoryID) => videoCategoryOptions.value.find((item) => Number(item.ID) === Number(categoryID))
  const getSeriesByID = (seriesID) => seriesOptions.value.find((item) => Number(item.ID) === Number(seriesID))
  const getEpisodeByID = (episodeID) => episodeOptions.value.find((item) => Number(item.ID) === Number(episodeID))

  const resolveVideoFolderBySeriesID = (seriesID) => {
    const series = getSeriesByID(seriesID)
    const category = getCategoryByID(series?.categoryId)
    const storageKey = normalizeStorageKey(category?.storageKey || '')
    return storageKey ? `english-learn/video/${storageKey}` : 'english-learn/video/general'
  }

  const episodeUploadFolder = computed(() => {
    if (episodeForm.value.seriesId) {
      return resolveVideoFolderBySeriesID(episodeForm.value.seriesId)
    }
    return 'english-learn/video/general'
  })

  const subtitleUploadFolder = computed(() => {
    const episode = getEpisodeByID(subtitleForm.value.episodeId)
    if (!episode) {
      return 'english-learn/video/general/subtitle'
    }
    return `${resolveVideoFolderBySeriesID(episode.seriesId)}/subtitle`
  })

  const resourceTypeLabel = (type) => {
    if (type === 'english_category') return '分类'
    if (type === 'video_series') return '剧集'
    if (type === 'video_episode') return '单集'
    return type || '-'
  }

  const loadManagedLanguages = async () => {
    try {
      const res = await getLanguageList({ page: 1, pageSize: 500 })
      const list = Array.isArray(res?.data) ? res.data : (res?.data?.list || [])
      managedLanguages.value = list
        .map((lang) => ({
          code: String(lang?.code || '').trim(),
          name: String(lang?.name || '').trim(),
          nativeName: String(lang?.nativeName || '').trim()
        }))
        .filter((lang) => lang.code)
    } catch (error) {
      managedLanguages.value = [
        { code: 'zh', name: '中文' },
        { code: 'mn', name: '蒙文' },
        { code: 'en', name: '英文' }
      ]
    }

    if (!managedLanguages.value.some((lang) => lang.code === displayLang.value)) {
      displayLang.value = managedLanguages.value[0]?.code || 'zh'
    }
    if (!managedLanguages.value.some((lang) => lang.code === sentencePreviewLang.value)) {
      sentencePreviewLang.value = managedLanguages.value[0]?.code || 'zh'
    }
  }

  const loadVideoCategoryList = async () => {
    const res = await getVideoCategoryList({ ...videoCategoryQuery.value })
    if (res.code !== 0) return
    videoCategoryTable.value = res.data?.list || []
    videoCategoryTotal.value = Number(res.data?.total || 0)
    videoCategoryQuery.value.page = Number(res.data?.page || videoCategoryQuery.value.page)
    videoCategoryQuery.value.pageSize = Number(res.data?.pageSize || videoCategoryQuery.value.pageSize)
  }

  const loadSeriesList = async () => {
    const params = { ...seriesQuery.value }
    if (!String(params.keyword || '').trim()) {
      delete params.keyword
    } else {
      params.keyword = String(params.keyword).trim()
    }
    const res = await getVideoSeriesList(params)
    if (res.code !== 0) return
    seriesTable.value = res.data?.list || []
    seriesTotal.value = Number(res.data?.total || 0)
    seriesQuery.value.page = Number(res.data?.page || seriesQuery.value.page)
    seriesQuery.value.pageSize = Number(res.data?.pageSize || seriesQuery.value.pageSize)
  }

  const loadEpisodeList = async () => {
    const params = { ...episodeQuery.value }
    if (!String(params.keyword || '').trim()) {
      delete params.keyword
    } else {
      params.keyword = String(params.keyword).trim()
    }
    const res = await getVideoEpisodeList(params)
    if (res.code !== 0) return
    episodeTable.value = res.data?.list || []
    episodeTotal.value = Number(res.data?.total || 0)
    episodeQuery.value.page = Number(res.data?.page || episodeQuery.value.page)
    episodeQuery.value.pageSize = Number(res.data?.pageSize || episodeQuery.value.pageSize)
  }

  const loadEntitlementList = async () => {
    const params = {
      page: entitlementQuery.value.page,
      pageSize: entitlementQuery.value.pageSize
    }
    if (entitlementQuery.value.userId) {
      params.userId = entitlementQuery.value.userId
    }
    if (entitlementQuery.value.resourceType) {
      params.resourceType = entitlementQuery.value.resourceType
    }

    const res = await getEntitlementList(params)
    if (res.code !== 0) return
    entitlementTable.value = res.data?.list || []
    entitlementTotal.value = Number(res.data?.total || 0)
    entitlementQuery.value.page = Number(res.data?.page || entitlementQuery.value.page)
    entitlementQuery.value.pageSize = Number(res.data?.pageSize || entitlementQuery.value.pageSize)
  }

  const loadVideoCategoryOptions = async () => {
    const res = await getVideoCategoryList({ page: 1, pageSize: 1000 })
    if (res.code !== 0) return
    videoCategoryOptions.value = res.data?.list || []
  }

  const loadSeriesOptions = async () => {
    const res = await getVideoSeriesList({ page: 1, pageSize: 1000 })
    if (res.code !== 0) return
    seriesOptions.value = res.data?.list || []
  }

  const loadEpisodeOptions = async () => {
    const res = await getVideoEpisodeList({ page: 1, pageSize: 1000 })
    if (res.code !== 0) return
    episodeOptions.value = res.data?.list || []
  }

  const loadVideoTagOptions = async () => {
    try {
      const res = await getVideoTagList({ page: 1, pageSize: 1000 })
      if (res.code === 0) {
        videoTagOptions.value = res.data?.list || []
      }
    } catch (e) {
      // ignore
    }
  }

  const handleVideoCategoryPageChange = (page) => {
    videoCategoryQuery.value.page = page
    loadVideoCategoryList()
  }

  const handleVideoCategorySizeChange = (size) => {
    videoCategoryQuery.value.pageSize = size
    videoCategoryQuery.value.page = 1
    loadVideoCategoryList()
  }

  const handleSeriesSearch = () => {
    seriesQuery.value.page = 1
    loadSeriesList()
  }

  const resetSeriesSearch = () => {
    seriesQuery.value = { page: 1, pageSize: 10, categoryId: undefined, keyword: '' }
    loadSeriesList()
  }

  const handleSeriesPageChange = (page) => {
    seriesQuery.value.page = page
    loadSeriesList()
  }

  const handleSeriesSizeChange = (size) => {
    seriesQuery.value.pageSize = size
    seriesQuery.value.page = 1
    loadSeriesList()
  }

  const handleEpisodeSearch = () => {
    episodeQuery.value.page = 1
    loadEpisodeList()
  }

  const resetEpisodeSearch = () => {
    episodeQuery.value = { page: 1, pageSize: 10, seriesId: undefined, keyword: '' }
    loadEpisodeList()
  }

  const handleEpisodePageChange = (page) => {
    episodeQuery.value.page = page
    loadEpisodeList()
  }

  const handleEpisodeSizeChange = (size) => {
    episodeQuery.value.pageSize = size
    episodeQuery.value.page = 1
    loadEpisodeList()
  }

  const handleEntitlementSearch = () => {
    entitlementQuery.value.page = 1
    loadEntitlementList()
  }

  const resetEntitlementSearch = () => {
    entitlementQuery.value = { page: 1, pageSize: 10, userId: undefined, resourceType: undefined }
    loadEntitlementList()
  }

  const handleEntitlementPageChange = (page) => {
    entitlementQuery.value.page = page
    loadEntitlementList()
  }

  const handleEntitlementSizeChange = (size) => {
    entitlementQuery.value.pageSize = size
    entitlementQuery.value.page = 1
    loadEntitlementList()
  }

  const openVideoCategoryDialog = (row) => {
    videoCategoryDialogMode.value = row?.ID ? 'edit' : 'create'
    videoCategoryForm.value = {
      ID: row?.ID || 0,
      nameI18n: normalizeI18nObject(row?.name || ''),
      storageKey: row?.storageKey || '',
      showHome: row?.showHome !== false,
      sort: Number(row?.sort || 0)
    }
    videoCategoryDialogVisible.value = true
  }

  const submitVideoCategory = async () => {
    const payload = {
      ID: videoCategoryForm.value.ID,
      name: stringifyI18nObject(videoCategoryForm.value.nameI18n),
      storageKey: normalizeStorageKey(videoCategoryForm.value.storageKey),
      showHome: videoCategoryForm.value.showHome !== false,
      sort: Number(videoCategoryForm.value.sort || 0)
    }

    if (!formatI18nText(payload.name)) {
      ElMessage.warning('分类名称不能为空')
      return
    }

    const res = videoCategoryDialogMode.value === 'create'
      ? await createVideoCategory(payload)
      : await updateVideoCategory(payload)

    if (res.code !== 0) return
    ElMessage.success(videoCategoryDialogMode.value === 'create' ? '视频分类创建成功' : '视频分类更新成功')
    videoCategoryDialogVisible.value = false
    await Promise.all([loadVideoCategoryList(), loadVideoCategoryOptions()])
  }

  const removeVideoCategory = (row) => {
    ElMessageBox.confirm(`确认删除视频分类【${formatI18nText(row.name)}】吗？`, '删除确认', {
      type: 'warning'
    }).then(async () => {
      const res = await deleteVideoCategory({ ID: row.ID })
      if (res.code !== 0) return
      ElMessage.success('视频分类删除成功')
      await Promise.all([
        loadVideoCategoryList(),
        loadVideoCategoryOptions(),
        loadSeriesList(),
        loadSeriesOptions(),
        loadEpisodeList(),
        loadEpisodeOptions()
      ])
    })
  }

  const openSeriesDialog = (row) => {
    seriesDialogMode.value = row?.ID ? 'edit' : 'create'
    seriesForm.value = {
      ID: row?.ID || 0,
      categoryId: row?.categoryId || seriesQuery.value.categoryId || undefined,
      nameI18n: normalizeI18nObject(row?.name || ''),
      coverId: Number(row?.coverId || 0),
      coverUrl: row?.coverUrl || '',
      price: Number(row?.price || 0),
      needVip: !!row?.needVip,
      showHome: row?.showHome !== false
    }
    seriesDialogVisible.value = true
  }

  const submitSeries = async () => {
    if (!seriesForm.value.categoryId) {
      ElMessage.warning('请选择所属分类')
      return
    }

    const payload = {
      ID: seriesForm.value.ID,
      categoryId: seriesForm.value.categoryId,
      name: stringifyI18nObject(seriesForm.value.nameI18n),
      coverId: Number(seriesForm.value.coverId || 0),
      coverUrl: String(seriesForm.value.coverUrl || '').trim(),
      price: Number(seriesForm.value.price || 0),
      needVip: !!seriesForm.value.needVip,
      showHome: seriesForm.value.showHome !== false
    }

    if (!formatI18nText(payload.name)) {
      ElMessage.warning('剧集名称不能为空')
      return
    }

    const res = seriesDialogMode.value === 'create'
      ? await createVideoSeries(payload)
      : await updateVideoSeries(payload)

    if (res.code !== 0) return
    ElMessage.success(seriesDialogMode.value === 'create' ? '剧集创建成功' : '剧集更新成功')
    seriesDialogVisible.value = false
    await Promise.all([loadSeriesList(), loadSeriesOptions(), loadEpisodeOptions()])
  }

  const removeSeries = (row) => {
    ElMessageBox.confirm(`确认删除剧集【${formatI18nText(row.name)}】吗？`, '删除确认', {
      type: 'warning'
    }).then(async () => {
      const res = await deleteVideoSeries({ ID: row.ID })
      if (res.code !== 0) return
      ElMessage.success('剧集删除成功')
      await Promise.all([loadSeriesList(), loadSeriesOptions(), loadEpisodeList(), loadEpisodeOptions()])
    })
  }

  const openEpisodeDialog = async (row) => {
    episodeDialogMode.value = row?.ID ? 'edit' : 'create'
    episodeForm.value = {
      ID: row?.ID || 0,
      seriesId: row?.seriesId || episodeQuery.value.seriesId || undefined,
      nameI18n: normalizeI18nObject(row?.name || ''),
      videoUrl: row?.videoUrl || '',
      trialPercent: Number(row?.trialPercent || 8),
      sort: Number(row?.sort || 0),
      tagIds: []
    }
    // Reset HLS state
    hlsMode.value = false
    ffmpegReady.value = false
    ffmpegChecking.value = false
    hlsSlicing.value = false
    hlsProgress.value = ''
    hlsFolder.value = episodeUploadFolder.value
    hlsFile.value = null
    if (hlsFileInput.value) {
      hlsFileInput.value.value = ''
    }
    // If editing an m3u8 episode, auto-enable HLS mode
    if (row?.videoType === 'm3u8') {
      hlsMode.value = true
      onHlsModeChange(true)
    }
    // Load existing tags when editing
    if (row?.ID) {
      try {
        const detailRes = await findVideoEpisode({ ID: row.ID })
        if (detailRes.code === 0 && detailRes.data?.tags) {
          episodeForm.value.tagIds = detailRes.data.tags.map(t => t.ID)
        }
      } catch (e) {
        // ignore tag load failure
      }
    }
    episodeDialogVisible.value = true
  }

  const onHlsModeChange = async (value) => {
    if (!value) {
      ffmpegReady.value = false
      ffmpegChecking.value = false
      hlsProgress.value = ''
      return
    }
    ffmpegChecking.value = true
    ffmpegReady.value = false
    try {
      const res = await checkFfmpeg()
      if (res.code === 0 && res.data?.available) {
        ffmpegReady.value = true
      } else {
        ffmpegReady.value = false
        ElMessage.warning('FFmpeg 不可用，无法使用 HLS 切片功能')
      }
    } catch (e) {
      ffmpegReady.value = false
      ElMessage.warning('FFmpeg 检测失败，请确认服务端已安装 FFmpeg')
    } finally {
      ffmpegChecking.value = false
    }
  }

  const onHlsFileChange = (e) => {
    hlsFile.value = e.target.files[0] || null
  }

  const onHlsSliceClick = async () => {
    if (!hlsFile.value) {
      ElMessage.warning('请选择视频文件')
      return
    }
    if (!episodeForm.value.seriesId) {
      ElMessage.warning('请先选择所属剧集')
      return
    }
    if (!formatI18nText(stringifyI18nObject(episodeForm.value.nameI18n))) {
      ElMessage.warning('单集名称不能为空')
      return
    }

    hlsSlicing.value = true
    hlsProgress.value = '正在切片上传中，请稍候（大视频可能需要几分钟）...'

    const formData = new FormData()
    formData.append('file', hlsFile.value)
    formData.append('seriesId', episodeForm.value.seriesId)
    formData.append('name', stringifyI18nObject(episodeForm.value.nameI18n))
    formData.append('trialPercent', episodeForm.value.trialPercent)
    formData.append('sort', episodeForm.value.sort)
    if (hlsFolder.value) {
      formData.append('folder', hlsFolder.value)
    }
    if (episodeForm.value.ID > 0) {
      formData.append('episodeId', episodeForm.value.ID)
    }

    try {
      const res = await sliceVideoEpisode(formData)
      if (res.code !== 0) {
        ElMessage.error('HLS 切片失败: ' + (res.msg || '未知错误'))
        return
      }
      ElMessage.success(res.msg || 'HLS 切片上传成功')
      if (res.data?.videoUrl) {
        episodeForm.value.videoUrl = res.data.videoUrl
      }
      episodeDialogVisible.value = false
      await Promise.all([loadEpisodeList(), loadEpisodeOptions()])
    } catch (e) {
      ElMessage.error('HLS 切片请求失败，请检查网络或服务端状态')
    } finally {
      hlsSlicing.value = false
      hlsProgress.value = ''
    }
  }

  const submitEpisode = async () => {
    if (!episodeForm.value.seriesId) {
      ElMessage.warning('请选择所属剧集')
      return
    }
    // In HLS mode, require video to be sliced first
    if (hlsMode.value && !String(episodeForm.value.videoUrl || '').trim()) {
      ElMessage.warning('请先选择视频文件进行切片上传')
      return
    }
    if (!String(episodeForm.value.videoUrl || '').trim()) {
      ElMessage.warning('视频地址不能为空')
      return
    }

    const payload = {
      ID: episodeForm.value.ID,
      seriesId: episodeForm.value.seriesId,
      name: stringifyI18nObject(episodeForm.value.nameI18n),
      videoUrl: String(episodeForm.value.videoUrl || '').trim(),
      trialPercent: Number(episodeForm.value.trialPercent || 8),
      sort: Number(episodeForm.value.sort || 0),
      tagIds: episodeForm.value.tagIds || []
    }

    if (!formatI18nText(payload.name)) {
      ElMessage.warning('单集名称不能为空')
      return
    }

    const res = episodeDialogMode.value === 'create'
      ? await createVideoEpisode(payload)
      : await updateVideoEpisode(payload)

    if (res.code !== 0) return
    ElMessage.success(episodeDialogMode.value === 'create' ? '单集创建成功' : '单集更新成功')
    episodeDialogVisible.value = false
    await Promise.all([loadEpisodeList(), loadEpisodeOptions()])
  }

  const removeEpisode = (row) => {
    ElMessageBox.confirm(`确认删除单集【${formatI18nText(row.name)}】吗？`, '删除确认', {
      type: 'warning'
    }).then(async () => {
      const res = await deleteVideoEpisode({ ID: row.ID })
      if (res.code !== 0) return
      ElMessage.success('单集删除成功')
      await Promise.all([loadEpisodeList(), loadEpisodeOptions()])
    })
  }

  const openSubtitleParser = async (row) => {
    const episodeId = Number(row?.ID || row?.id || 0)

    // 先重置表单
    subtitleForm.value = {
      episodeId,
      englishSubtitleUrl: '',
      translationSubtitleMap: {}
    }
    activeTab.value = 'episode'

    // 自动回填已上传的字幕文件路径（从 VideoSubtitle 表查询）
    try {
      const subtitleRes = await getEpisodeSubtitles({ episodeId })
      if (subtitleRes.code === 0 && Array.isArray(subtitleRes.data) && subtitleRes.data.length > 0) {
        const subtitleList = subtitleRes.data
        const newMap = {}
        let englishUrl = ''
        for (const sub of subtitleList) {
          const lang = String(sub.language || '').trim()
          const url = String(sub.subtitleUrl || '').trim()
          if (!lang || !url) continue
          if (lang === 'en') {
            englishUrl = url
          } else {
            newMap[lang] = url
          }
        }
        // 整体替换以触发 Vue 响应式更新
        subtitleForm.value = {
          episodeId,
          englishSubtitleUrl: englishUrl,
          translationSubtitleMap: newMap
        }
        const filledCount = (englishUrl ? 1 : 0) + Object.keys(newMap).length
        ElMessage.success(`已自动回填 ${filledCount} 个字幕文件路径`)
      }
    } catch (e) {
      console.error('回填字幕文件路径失败:', e)
    }

    // 预检该单集是否已解析过字幕句子
    try {
      const res = await getVideoSentenceList({ episodeId })
      if (res.code === 0 && Array.isArray(res.data) && res.data.length > 0) {
        ElMessage.info('该单集已存在解析后的字幕句子，可直接点击"查看字幕"进行预览或编辑')
      }
    } catch (e) {
      // ignore 预检失败不影响后续解析操作
    }

    if (subtitleSectionRef.value?.scrollIntoView) {
      setTimeout(() => {
        subtitleSectionRef.value.scrollIntoView({ behavior: 'smooth', block: 'start' })
      }, 30)
    }
  }

  const previewEpisodeSubtitles = async (row) => {
    const episodeId = Number(row?.ID || row?.id || 0)
    if (!episodeId) {
      ElMessage.warning('单集ID无效')
      return
    }

    const res = await getVideoSentenceList({ episodeId })
    if (res.code !== 0) {
      return
    }

    sentencePreviewEpisodeId.value = episodeId
    sentencePreviewEpisodeName.value = formatI18nText(row?.name)
    sentencePreviewTable.value = (Array.isArray(res.data) ? res.data : []).map((item) => {
      const translateObj = parseI18nObject(item?.translate)
      const lang = String(sentencePreviewLang.value || 'zh')
      return {
        ...item,
        translateObj,
        translateText: translateObj[lang] || translateObj.zh || translateObj.en || translateObj.mn || ''
      }
    })
    sentencePreviewOriginalMap.value = buildSentencePreviewSnapshotMap(sentencePreviewTable.value)
    sentencePreviewTab.value = 'sentences'
    sentencePreviewVisible.value = true

    // 预加载关键词列表
    episodeKeywordList.value = []
    previewSelectedKeywordIds.value = []
    loadEpisodeKeywords()
  }

  const syncSentenceTranslate = (row) => {
    if (!row || typeof row !== 'object') return
    const lang = String(sentencePreviewLang.value || 'zh')
    const next = { ...(row.translateObj || {}) }
    next[lang] = String(row.translateText || '')
    row.translateObj = next
    row.translate = JSON.stringify(next)
  }

  const refreshSentenceTranslateByLang = () => {
    const lang = String(sentencePreviewLang.value || 'zh')
    sentencePreviewTable.value = sentencePreviewTable.value.map((item) => {
      const translateObj = parseI18nObject(item?.translateObj || item?.translate)
      return {
        ...item,
        translateObj,
        translateText: translateObj[lang] || translateObj.zh || translateObj.en || translateObj.mn || ''
      }
    })
  }

  const saveSentencePreview = async () => {
    if (!sentencePreviewEpisodeId.value) {
      ElMessage.warning('单集ID无效')
      return
    }
    sentenceSaving.value = true
    const payload = {
      episodeId: sentencePreviewEpisodeId.value,
      sentences: pendingSentenceChanges.value.map((item) => ({
        id: Number(item?.id || item?.ID || 0),
        english: String(item?.english || '').trim(),
        translate: JSON.stringify(item?.translateObj || {}),
        startTime: Number(item?.startTime || 0),
        endTime: Number(item?.endTime || 0)
      }))
    }
    if (payload.sentences.length === 0) {
      sentenceSaving.value = false
      ElMessage.info('没有可保存的改动')
      return
    }
    const res = await updateVideoSentenceList(payload)
    sentenceSaving.value = false
    if (res.code !== 0) return
    ElMessage.success('字幕修改已保存')
    sentencePreviewOriginalMap.value = buildSentencePreviewSnapshotMap(sentencePreviewTable.value)
    sentencePreviewVisible.value = false
  }

  const resetSentencePreviewChanges = () => {
    const lang = String(sentencePreviewLang.value || 'zh')
    sentencePreviewTable.value = sentencePreviewTable.value.map((row) => {
      const rowId = Number(row?.id || row?.ID || 0)
      const original = sentencePreviewOriginalMap.value[rowId]
      if (!original) return row
      const translateObj = { ...original.translateObj }
      return {
        ...row,
        english: original.english,
        startTime: original.startTime,
        endTime: original.endTime,
        translateObj,
        translate: JSON.stringify(translateObj),
        translateText: translateObj[lang] || translateObj.zh || translateObj.en || translateObj.mn || ''
      }
    })
  }

  const closeSentencePreview = () => {
    if (pendingSentenceChanges.value.length === 0) {
      sentencePreviewVisible.value = false
      return
    }
    ElMessageBox.confirm(`当前有 ${pendingSentenceChanges.value.length} 条未保存改动，确认关闭？`, '未保存改动', {
      type: 'warning'
    }).then(() => {
      sentencePreviewVisible.value = false
    })
  }

  const handleSentencePreviewBeforeClose = (done) => {
    if (pendingSentenceChanges.value.length === 0) {
      done()
      return
    }
    ElMessageBox.confirm(`当前有 ${pendingSentenceChanges.value.length} 条未保存改动，确认关闭？`, '未保存改动', {
      type: 'warning'
    }).then(() => done())
  }

  watch(
    () => sentencePreviewLang.value,
    () => {
      refreshSentenceTranslateByLang()
    }
  )

  // 切换到"全部单词"tab 时自动选中已在词库的单词
  watch(
    () => sentencePreviewTab.value,
    (newTab) => {
      if (newTab !== 'allWords') return
      autoSelectAllMatchedWords()
    }
  )
  // 全部单词列表变化时（如关键词数据加载完成）也触发自动选中
  watch(
    () => episodeKeywordList.value,
    () => {
      if (sentencePreviewTab.value !== 'allWords') return
      autoSelectAllMatchedWords()
    }
  )

  const autoSelectAllMatchedWords = () => {
    if (!allWordsTableRef.value) return
    setTimeout(() => {
      if (!allWordsTableRef.value) return
      allWordsTableRef.value.clearSelection()
      const matchedRows = allSubtitleWords.value.filter(k => k.matched)
      matchedRows.forEach(row => {
        allWordsTableRef.value.toggleRowSelection(row, true)
      })
    }, 60)
  }

  // 关键词扫描与解析（两步流程）
  const handleScanKeywords = async () => {
    if (!subtitleForm.value.episodeId) {
      ElMessage.warning('请选择目标单集')
      return
    }

    const englishSubtitleURL = String(subtitleForm.value.englishSubtitleUrl || '').trim()
    if (!englishSubtitleURL) {
      ElMessage.warning('英文字幕文件不能为空')
      return
    }

    keywordScanning.value = true
    keywordScanDone.value = false
    keywordScanResults.value = []
    selectedKeywordIds.value = []

    const res = await scanKeywords({
      episodeId: subtitleForm.value.episodeId,
      englishSubtitleUrl: englishSubtitleURL
    })

    keywordScanning.value = false
    if (res.code !== 0) return

    keywordScanResults.value = res.data || []
    keywordScanDone.value = true

    // 默认选中所有已匹配的单词
    const matchedIds = keywordScanResults.value.filter(k => k.matched).map(k => k.wordId)
    selectedKeywordIds.value = [...matchedIds]

    // 同步勾选表格
    setTimeout(() => {
      if (keywordTableRef.value) {
        const matchedRows = keywordScanResults.value.filter(k => k.matched)
        matchedRows.forEach(row => {
          keywordTableRef.value.toggleRowSelection(row, true)
        })
      }
    }, 50)

    ElMessage.success(`扫描完成，共提取 ${keywordScanResults.value.length} 个单词，其中 ${matchedIds.length} 个已匹配词库`)
  }

  const handleConfirmParse = async () => {
    if (!subtitleForm.value.episodeId) {
      ElMessage.warning('请选择目标单集')
      return
    }

    const englishSubtitleURL = String(subtitleForm.value.englishSubtitleUrl || '').trim()
    if (!englishSubtitleURL) {
      ElMessage.warning('英文字幕文件不能为空')
      return
    }

    if (selectedKeywordIds.value.length === 0) {
      ElMessage.warning('请至少选择一个重点单词')
      return
    }

    const translationSubtitle = Object.entries(subtitleForm.value.translationSubtitleMap || {})
      .map(([language, subtitleUrl]) => ({
        language,
        subtitleUrl: String(subtitleUrl || '').trim()
      }))
      .filter((item) => item.language && item.subtitleUrl)

    const res = await parseSubtitleFiles({
      episodeId: subtitleForm.value.episodeId,
      englishSubtitleUrl: englishSubtitleURL,
      translationSubtitle,
      keywordIds: selectedKeywordIds.value
    })

    if (res.code !== 0) return
    ElMessage.success('字幕文件解析并入库成功')

    // 清除扫描状态
    keywordScanDone.value = false
    keywordScanResults.value = []
    selectedKeywordIds.value = []
  }

  const handleKeywordSelectionChange = (rows) => {
    selectedKeywordIds.value = rows.map(r => r.wordId)
  }

  const selectAllKeywords = (selectAll) => {
    if (!keywordTableRef.value) return
    if (selectAll) {
      const matchedRows = keywordScanResults.value.filter(k => k.matched)
      matchedRows.forEach(row => keywordTableRef.value.toggleRowSelection(row, true))
    } else {
      keywordTableRef.value.clearSelection()
    }
  }

  // 预览弹窗-关键词管理
  const loadEpisodeKeywords = async () => {
    if (!sentencePreviewEpisodeId.value) return
    const res = await getEpisodeKeywords({ episodeId: sentencePreviewEpisodeId.value })
    if (res.code !== 0) return
    episodeKeywordList.value = res.data || []

    // 默认选中所有已匹配的
    const matchedIds = episodeKeywordList.value.filter(k => k.matched).map(k => k.wordId)
    previewSelectedKeywordIds.value = [...matchedIds]

    setTimeout(() => {
      if (previewKeywordTableRef.value) {
        previewKeywordTableRef.value.clearSelection()
        const matchedRows = episodeKeywordList.value.filter(k => k.matched)
        matchedRows.forEach(row => {
          previewKeywordTableRef.value.toggleRowSelection(row, true)
        })
      }
    }, 50)
  }

  const handlePreviewKeywordSelectionChange = (rows) => {
    previewSelectedKeywordIds.value = rows.map(r => r.wordId)
  }

  const selectAllPreviewKeywords = (selectAll) => {
    if (!previewKeywordTableRef.value) return
    if (selectAll) {
      const matchedRows = episodeKeywordList.value.filter(k => k.matched)
      matchedRows.forEach(row => previewKeywordTableRef.value.toggleRowSelection(row, true))
    } else {
      previewKeywordTableRef.value.clearSelection()
    }
  }

  const handleRehighlightPreview = async () => {
    if (!sentencePreviewEpisodeId.value) {
      ElMessage.warning('单集ID无效')
      return
    }
    if (previewSelectedKeywordIds.value.length === 0) {
      ElMessage.warning('请至少选择一个重点单词')
      return
    }

    previewRehighlighting.value = true
    const res = await rehighlightSentences({
      episodeId: sentencePreviewEpisodeId.value,
      keywordIds: previewSelectedKeywordIds.value
    })
    previewRehighlighting.value = false

    if (res.code !== 0) return
    ElMessage.success('字幕重新高亮成功')

    // 重新加载句子列表以显示新的高亮结果
    const sentenceRes = await getVideoSentenceList({ episodeId: sentencePreviewEpisodeId.value })
    if (sentenceRes.code === 0) {
      const lang = String(sentencePreviewLang.value || 'zh')
      sentencePreviewTable.value = (Array.isArray(sentenceRes.data) ? sentenceRes.data : []).map((item) => {
        const translateObj = parseI18nObject(item?.translate)
        return {
          ...item,
          translateObj,
          translateText: translateObj[lang] || translateObj.zh || translateObj.en || translateObj.mn || ''
        }
      })
      sentencePreviewOriginalMap.value = buildSentencePreviewSnapshotMap(sentencePreviewTable.value)
    }

    // 刷新关键词列表
    await loadEpisodeKeywords()
  }

  // 全部单词 tab - 选择管理
  const handleAllWordsSelectionChange = (rows) => {
    allWordsSelectedIds.value = rows.map(r => r.wordId).filter(id => id > 0)
  }

  const selectAllAllWords = (selectAll) => {
    if (!allWordsTableRef.value) return
    if (selectAll) {
      const matchedRows = allSubtitleWords.value.filter(k => k.matched)
      matchedRows.forEach(row => allWordsTableRef.value.toggleRowSelection(row, true))
    } else {
      allWordsTableRef.value.clearSelection()
    }
  }

  const handleRehighlightAllWords = async () => {
    if (!sentencePreviewEpisodeId.value) {
      ElMessage.warning('单集ID无效')
      return
    }
    if (allWordsSelectedIds.value.length === 0) {
      ElMessage.warning('请至少选择一个已入库的单词')
      return
    }

    allWordsRehighlighting.value = true

    // 收集所有选中行（包括未入库的），找出未入库的单词
    const selectedRows = allWordsTableRef.value?.getSelectionRows?.() || []
    const unmatchedWords = selectedRows
      .filter(r => !r.matched)
      .map(r => r.word)

    let newIds = []
    if (unmatchedWords.length > 0) {
      // 先批量入库未匹配的单词
      const batchRes = await batchCreateWords({ words: unmatchedWords })
      if (batchRes.code === 0 && Array.isArray(batchRes.data)) {
        newIds = batchRes.data.map(item => item.wordId)
        ElMessage.success(`已自动入库 ${newIds.length} 个新单词`)
      }
    }

    // 合并已有 wordId + 新入库的 wordId
    const allIds = [...new Set([...allWordsSelectedIds.value, ...newIds])]
    const res = await rehighlightSentences({
      episodeId: sentencePreviewEpisodeId.value,
      keywordIds: allIds
    })
    allWordsRehighlighting.value = false

    if (res.code !== 0) return
    ElMessage.success(`字幕重新高亮成功，已应用 ${allIds.length} 个重点单词`)

    // 重新加载句子列表以显示新的高亮结果
    const sentenceRes = await getVideoSentenceList({ episodeId: sentencePreviewEpisodeId.value })
    if (sentenceRes.code === 0) {
      const lang = String(sentencePreviewLang.value || 'zh')
      sentencePreviewTable.value = (Array.isArray(sentenceRes.data) ? sentenceRes.data : []).map((item) => {
        const translateObj = parseI18nObject(item?.translate)
        return {
          ...item,
          translateObj,
          translateText: translateObj[lang] || translateObj.zh || translateObj.en || translateObj.mn || ''
        }
      })
      sentencePreviewOriginalMap.value = buildSentencePreviewSnapshotMap(sentencePreviewTable.value)
    }

    // 刷新关键词列表
    await loadEpisodeKeywords()
    // 清空全部单词选择
    allWordsSelectedIds.value = []
    if (allWordsTableRef.value) {
      allWordsTableRef.value.clearSelection()
    }
  }

  const openEntitlementDialog = () => {
    entitlementForm.value = {
      userId: undefined,
      resourceType: 'video_series',
      resourceId: undefined,
      expireAt: '',
      remark: ''
    }
    entitlementDialogVisible.value = true
  }

  const submitEntitlement = async () => {
    if (!entitlementForm.value.userId) {
      ElMessage.warning('用户ID不能为空')
      return
    }
    if (!entitlementForm.value.resourceType) {
      ElMessage.warning('资源类型不能为空')
      return
    }
    if (!entitlementForm.value.resourceId) {
      ElMessage.warning('资源ID不能为空')
      return
    }

    const payload = {
      userId: Number(entitlementForm.value.userId),
      resourceType: entitlementForm.value.resourceType,
      resourceId: Number(entitlementForm.value.resourceId),
      remark: entitlementForm.value.remark || ''
    }
    if (entitlementForm.value.expireAt) {
      payload.expireAt = entitlementForm.value.expireAt
    }

    const res = await grantEntitlement(payload)
    if (res.code !== 0) return
    ElMessage.success('授权成功')
    entitlementDialogVisible.value = false
    loadEntitlementList()
  }

  const removeEntitlement = (row) => {
    ElMessageBox.confirm(`确认撤销用户 ${row.userId} 的授权记录 #${row.ID} 吗？`, '撤销确认', {
      type: 'warning'
    }).then(async () => {
      const res = await revokeEntitlement({
        userId: row.userId,
        resourceType: row.resourceType,
        resourceId: row.resourceId
      })
      if (res.code !== 0) return
      ElMessage.success('授权已撤销')
      loadEntitlementList()
    })
  }

  onMounted(async () => {
    await Promise.all([
      loadVideoCategoryList(),
      loadSeriesList(),
      loadEpisodeList(),
      loadEntitlementList(),
      loadVideoCategoryOptions(),
      loadSeriesOptions(),
      loadEpisodeOptions(),
      loadManagedLanguages(),
      loadVideoTagOptions()
    ])
  })
</script>

<style scoped>
  .lang-switch-row {
    margin-top: 12px;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .subtitle-preview-toolbar {
    margin-bottom: 12px;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .toolbar-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    margin-bottom: 16px;
    flex-wrap: wrap;
  }

  .toolbar-tip {
    color: #4b5563;
    font-size: 13px;
  }

  .subtitle-box {
    margin-top: 12px;
    padding: 12px 0;
  }

  .upload-inline {
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    gap: 10px;
    width: 100%;
  }

  .upload-inline-wide {
    max-width: 680px;
  }

  .dialog-hint {
    margin-top: 6px;
    color: #6b7280;
    font-size: 12px;
    line-height: 1.4;
  }

  .mt-2 {
    margin-top: 8px;
  }

  .w-full {
    width: 100%;
  }

  .keyword-scan-result {
    margin-top: 16px;
    padding: 12px 16px;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    background: #fafafa;
  }

  .keyword-toolbar {
    margin-bottom: 10px;
    display: flex;
    gap: 8px;
  }

  .hls-folder-bar {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 10px;
    font-size: 13px;
  }

  .hls-folder-label {
    white-space: nowrap;
    color: #909399;
  }

  .hls-folder-input {
    width: 280px;
  }
</style>
