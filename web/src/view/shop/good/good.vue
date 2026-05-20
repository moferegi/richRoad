<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        :rules="searchRule"
        @keyup.enter="onSubmit"
      >
        <el-form-item
          label="创建日期"
          prop="createdAt"
        >
          <template #label>
            <span>
              创建日期
              <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker
            v-model="searchInfo.startCreatedAt"
            type="datetime"
            placeholder="开始日期"
            :disabled-date="time=> searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"
          />
          —
          <el-date-picker
            v-model="searchInfo.endCreatedAt"
            type="datetime"
            placeholder="结束日期"
            :disabled-date="time=> searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"
          />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-select v-model="searchInfo.status" clearable placeholder="请选择">
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>

        <el-form-item label="预售">
          <el-select v-model="searchInfo.isPresale" clearable placeholder="全部" style="width:100px;">
            <el-option label="预售" :value="true" />
            <el-option label="非预售" :value="false" />
          </el-select>
        </el-form-item>

        <el-form-item label="关键词">
          <el-input v-model="searchInfo.keyword" placeholder="名称/描述/标签" clearable style="width:160px;" />
        </el-form-item>

        <el-form-item label="商品分类">
          <el-select v-model="searchInfo.categoryID" clearable placeholder="全部分类" style="width:160px;">
            <el-option v-for="cat in categoryList" :key="cat.ID" :label="cat.title" :value="cat.ID" />
          </el-select>
        </el-form-item>

        <el-form-item label="排序">
          <el-select v-model="searchInfo.orderBy" clearable placeholder="排序字段" style="width:120px;">
            <el-option label="浏览量" value="view_num" />
            <el-option label="销量" value="sale_num" />
            <el-option label="收藏数" value="collect_num" />
            <el-option label="价格" value="price" />
            <el-option label="创建时间" value="created_at" />
          </el-select>
          <el-select v-model="searchInfo.orderDir" clearable placeholder="方向" style="width:90px; margin-left:4px;">
            <el-option label="升序" value="asc" />
            <el-option label="降序" value="desc" />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            icon="search"
            @click="onSubmit"
          >查询</el-button>
          <el-button
            icon="refresh"
            @click="onReset"
          >重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          type="primary"
          icon="plus"
          @click="openDialog"
        >新增</el-button>
        <el-button
          icon="delete"
          style="margin-left: 10px;"
          :disabled="!multipleSelection.length"
          @click="onDelete"
        >删除</el-button>
      </div>
      <div v-if="enabledLangs.length" style="display:flex;align-items:center;justify-content:flex-end;margin-bottom:8px;">
        <span style="margin-right:6px;font-size:13px;color:#666;">表格语言：</span>
        <el-select v-model="tableLang" size="small" style="width:120px;">
          <el-option label="默认" value="" />
          <el-option v-for="l in enabledLangs" :key="l.code" :label="l.name || l.code" :value="l.code" />
        </el-select>
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        :default-sort="{ prop: 'ID', order: 'descending' }"
        max-height="70vh"
        @selection-change="handleSelectionChange"
      >
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column align="left" label="ID" prop="ID" width="70" sortable />

        <el-table-column
          align="left"
          label="日期"
          width="180"
          sortable
        >
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="商品名称"
          width="120"
        >
          <template #default="scope">{{ formatI18nField(scope.row.title) }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="商品描述"
          width="120"
        >
          <template #default="scope">{{ formatI18nField(scope.row.description) }}</template>
        </el-table-column>
        <el-table-column
          label="商品图片"
          width="200"
        >
          <template #default="scope">
            <div class="good-image-wrap">
              <el-image
                class="good-image"
                :src="getPrimaryImageLink(scope.row)"
                fit="cover"
              />
              <span v-if="getImageSizeText(scope.row)" class="good-image-size-badge">{{ getImageSizeText(scope.row) }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          label="上装图"
          width="120"
        >
          <template #default="scope">
            <el-image
              v-if="scope.row.upperImage"
              style="width: 60px; height: 60px"
              :src="getUrl(scope.row.upperImage)"
              fit="cover"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column
          label="下装图"
          width="120"
        >
          <template #default="scope">
            <el-image
              v-if="scope.row.lowerImage"
              style="width: 60px; height: 60px"
              :src="getUrl(scope.row.lowerImage)"
              fit="cover"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="商品价格"
          width="120"
        >
          <template #default="scope">{{ formatTablePrice(scope.row) }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="商品评分"
          prop="rating"
          width="120"
        />
        <el-table-column
          align="left"
          label="商品评论数量"
          prop="reviewCount"
          width="120"
        />
        <el-table-column
          align="left"
          label="浏览量"
          prop="view_num"
          width="90"
        />
        <el-table-column
          align="left"
          label="销量"
          prop="saleNum"
          width="90"
        />
        <el-table-column
          align="left"
          label="收藏数"
          prop="collect_num"
          width="90"
        />
        <el-table-column
          align="left"
          label="总库存"
          width="90"
        >
          <template #default="scope">{{ getTotalInventory(scope.row) }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="商品类型"
          prop="categoryID"
          width="120"
        >
          <template #default="scope">{{ categoryList.find(item => item.ID === scope.row.categoryID)?.title }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="邮费"
          prop="postage"
          width="120"
        />
        <el-table-column
          align="left"
          label="折扣"
          prop="discount"
          width="120"
        />

        <el-table-column
          align="left"
          label="推荐"
          prop="recommend"
          width="120"
        >
          <template #default="scope">{{ formatBoolean(scope.row.recommend) }}</template>
        </el-table-column>

        <el-table-column
          align="left"
          label="备注"
          prop="remark"
          width="120"
          show-overflow-tooltip
        />

        <el-table-column
          align="left"
          label="预售"
          width="80"
        >
          <template #default="scope">
            <el-tag v-if="scope.row.isPresale" type="warning" size="small">预售</el-tag>
          </template>
        </el-table-column>

        <el-table-column
          align="left"
          label="启用"
          prop="status"
          width="120"
        >
          <template #default="scope">
            <el-switch
              v-model="scope.row.status"
              inline-prompt
              :active-value="true"
              :inactive-value="false"
              @change="onStatusChange(scope.row)"
            />
          </template>
        </el-table-column>

        <el-table-column
          align="left"
          label="操作"
          fixed="right"
          min-width="240"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="edit"
              class="table-button"
              @click="updateGoodFunc(scope.row)"
            >变更</el-button>
            <el-button
              type="primary"
              link
              icon="InfoFilled"
              class="table-button"
              @click="setSKU(scope.row)"
            >SKU</el-button>
            <el-button
              type="primary"
              link
              icon="Goods"
              class="table-button"
              @click="gotoPurchase(scope.row)"
            >进货</el-button>
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteRow(scope.row)"
            >删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    <el-drawer
      v-model="dialogFormVisible"
      size="800"
      :show-close="false"
      :before-close="closeDialog"
      destroy-on-close
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type==='create'?'添加':'修改' }}</span>
          <div>
            <el-button
              type="primary"
              @click="enterDialog"
            >确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="elFormRef"
        :model="formData"
        label-position="top"
        :rules="rule"
        label-width="80px"
      >
        <el-row :gutter="12">
          <el-col :span="6">
            <el-form-item
              label="商品名称:"
              prop="title"
            >
              <el-input
                v-model="formData.title"
                :clearable="true"
                placeholder="默认商品名称"
              />
            </el-form-item>
            <MultiLangEditor
              :model="titleI18n"
              :languages="enabledLangs"
              title="商品名称多语言"
            />
          </el-col>
          <el-col :span="6">
            <el-form-item
              label="商品类型:"
              prop="categoryID"
            >
              <el-tree-select
                v-model="formData.categoryID"
                :data="categoryList"
                :props="{ label: 'title', value: 'ID', children: 'children', isLeaf: 'isLeaf' }"
                placeholder="请选择商品类型"
                check-strictly
                :only-leaf-select="true"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item
              label="状态:"
              prop="status"
            >
              <el-switch
                v-model="formData.status"
                active-color="#13ce66"
                inactive-color="#ff4949"
                active-text="是"
                inactive-text="否"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item
              label="推荐:"
              prop="recommend"
            >
              <el-switch
                v-model="formData.recommend"
                active-color="#13ce66"
                inactive-color="#ff4949"
                active-text="是"
                inactive-text="否"
                clearable
              />
            </el-form-item>
          </el-col>

        </el-row>
        <el-form-item
          label="商品图片(上传):"
          prop="imageUrl"
        >
          <SelectImage
            v-model="formData.imageUrl"
            file-type="image"
            :default-folder="GOOD_UPLOAD_FOLDER"
            :fixed-upload-folder="true"
          />
        </el-form-item>
        <el-form-item
          label="上装图(上传):"
          prop="upperImage"
        >
          <SelectImage
            v-model="formData.upperImage"
            file-type="image"
            :default-folder="GOOD_UPLOAD_FOLDER"
            :fixed-upload-folder="true"
          />
        </el-form-item>
        <el-form-item
          label="下装图(上传):"
          prop="lowerImage"
        >
          <SelectImage
            v-model="formData.lowerImage"
            file-type="image"
            :default-folder="GOOD_UPLOAD_FOLDER"
            :fixed-upload-folder="true"
          />
        </el-form-item>
        <el-form-item label="商品图片外链(优先于上传):" prop="externalImagePath">
          <el-input v-model="formData.externalImagePath" placeholder="相对路径如 /images/product.jpg 自动拼接外部域名" clearable />
          <div v-if="extDomain" class="text-xs text-gray-400 mt-1">当前外部域名: {{ extDomain }}</div>
        </el-form-item>
        <el-form-item
          label="商品轮播图:"
          prop="banner"
        >
          <div class="w-full">
            <div v-for="(item, index) in formData.banner" :key="index" class="mb-4 p-3" style="border: 1px solid #ebeef5; border-radius: 4px;">
              <el-row :gutter="12" class="items-center">
                <el-col :span="4">
                  <el-tag size="small">{{ index + 1 }}</el-tag>
                  <el-select v-model="item.type" placeholder="类型" size="small" style="width:80px; margin-left:4px;">
                    <el-option label="图片" value="image" />
                    <el-option label="视频" value="video" />
                  </el-select>
                </el-col>
                <el-col :span="16">
                  <el-form-item label="资源(上传):" class="mb-0" label-width="auto">
                    <SelectImage
                      v-model="item.url"
                      :file-type="item.type === 'video' ? 'video' : 'image'"
                      :default-folder="GOOD_BANNER_UPLOAD_FOLDER"
                      :fixed-upload-folder="true"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="4" class="text-right">
                  <el-button type="danger" icon="delete" size="small" @click="formData.banner.splice(index, 1)">删除</el-button>
                </el-col>
              </el-row>
              <el-form-item label="外部链接(优先于上传):" class="mt-2 mb-0" label-width="auto">
                <el-input v-model="item.externalUrl" placeholder="https://example.com/media.jpg" clearable />
              </el-form-item>
              <el-row :gutter="12" class="mt-2">
                <el-col :span="12">
                  <el-form-item label="文字(多语言):" class="mb-0" label-width="auto">
                    <el-input v-model="item.text" placeholder="默认文字" clearable />
                    <div class="w-full mt-2">
                      <MultiLangEditor
                        :model="item.textI18n"
                        :languages="enabledLangs"
                        :title="'轮播文案多语言 #' + (index + 1)"
                      />
                    </div>
                  </el-form-item>
                </el-col>
                <el-col :span="4">
                  <el-form-item label="颜色:" class="mb-0" label-width="auto">
                    <el-color-picker v-model="item.textColor" show-alpha size="small" />
                  </el-form-item>
                </el-col>
                <el-col :span="4">
                  <el-form-item label="字号:" class="mb-0" label-width="auto">
                    <el-input-number v-model="item.textSize" :min="8" :max="72" size="small" style="width:80px;" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="位置:" class="mb-0" label-width="auto">
                    <el-select v-model="item.textPosition" size="small">
                      <el-option label="顶部" value="top" />
                      <el-option label="居中" value="center" />
                      <el-option label="底部" value="bottom" />
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
            </div>
            <el-button type="primary" icon="plus" @click="formData.banner.push({ url: '', type: 'image', externalUrl: '', text: '', textI18n: {}, textColor: '#FFFFFF', textSize: 14, textPosition: 'bottom' })">添加轮播项</el-button>
          </div>
        </el-form-item>
        <el-form-item
          label="商品标签:"
          prop="tags"
        >
          <el-select
            v-model="formData.tags"
            placeholder="请选择商品标签"
            multiple
            clearable
            value-key="ID"
          >
            <el-option
              v-for="item in tags"
              :key="item.ID"
              :label="item.name"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          label="商品描述:"
          prop="description"
        >
          <el-input
            v-model="formData.description"
            :clearable="true"
            placeholder="默认商品描述"
          />
          <div class="w-full mt-2">
            <MultiLangEditor
              :model="descI18n"
              :languages="enabledLangs"
              title="商品描述多语言"
            />
          </div>
        </el-form-item>

        <el-row :gutter="12">
          <el-col :span="8">
            <el-form-item
              label="商品价格(单位：分):"
              prop="price"
            >
              <div class="w-full">
                <el-input-number
                  v-model="formData.price"
                  style="width:100%"
                  :precision="2"
                  :clearable="true"
                />
                <div class="mt-2 flex items-center justify-between gap-2">
                  <el-button type="primary" plain size="small" @click="openPriceRateDialog">汇率换算</el-button>
                  <span class="text-xs text-gray-500">多语言价格槽位：{{ Object.keys(priceI18n).length }}</span>
                </div>
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item
              label="邮费:"
              prop="postage"
            >
              <el-input-number
                v-model="formData.postage"
                style="width:100%"
                :precision="2"
                :clearable="true"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item
              label="折扣(%):"
              prop="discount"
            >
              <el-input-number
                v-model="formData.discount"
                :clearable="true"
                style="width:100%"
                :precision="0"
                placeholder="请输入折扣(%)"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="商品备注(来源等):" prop="remark">
          <el-input v-model="formData.remark" type="textarea" :rows="2" placeholder="请输入商品备注" />
        </el-form-item>

        <el-divider content-position="left">积分抵扣设置</el-divider>
        <el-row :gutter="12">
          <el-col :span="6">
            <el-form-item label="允许积分抵扣:" prop="pointsEnabled">
              <el-switch v-model="formData.pointsEnabled" />
            </el-form-item>
          </el-col>
          <el-col :span="9">
            <el-form-item label="最多可用积分数:" prop="pointsMaxUse">
              <el-input-number v-model="formData.pointsMaxUse" :min="0" style="width:100%;" :disabled="!formData.pointsEnabled" />
            </el-form-item>
          </el-col>
          <el-col :span="9">
            <el-form-item label="同商品可用积分次数:" prop="pointsUseTimes">
              <el-input-number v-model="formData.pointsUseTimes" :min="0" style="width:100%;" :disabled="!formData.pointsEnabled" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">预售设置</el-divider>
        <el-row :gutter="12">
          <el-col :span="6">
            <el-form-item label="开启预售:" prop="isPresale">
              <el-switch v-model="formData.isPresale" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="预售开关:" prop="presaleEnabled">
              <el-switch v-model="formData.presaleEnabled" :disabled="!formData.isPresale" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="预售数量:" prop="presaleQty">
              <el-input-number v-model="formData.presaleQty" :min="0" style="width:100%;" :disabled="!formData.isPresale" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="预售排序:" prop="presaleSort">
              <el-input-number v-model="formData.presaleSort" :min="0" style="width:100%;" :disabled="!formData.isPresale" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="预售开始:" prop="presaleStart">
              <el-date-picker
                v-model="formData.presaleStart"
                type="datetime"
                placeholder="选择预售开始时间"
                style="width:100%;"
                :disabled="!formData.isPresale"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="预售结束:" prop="presaleEnd">
              <el-date-picker
                v-model="formData.presaleEnd"
                type="datetime"
                placeholder="选择预售结束时间"
                style="width:100%;"
                :disabled="!formData.isPresale"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="预售弹窗开关:" prop="presalePopupEnabled" v-if="formData.isPresale">
          <el-switch v-model="formData.presalePopupEnabled" />
        </el-form-item>
        <el-form-item label="预售弹窗标题(多语言):" prop="presalePopupTitle" v-if="formData.isPresale && formData.presalePopupEnabled">
          <el-input v-model="formData.presalePopupTitle" placeholder="默认预售弹窗标题" class="mb-2" />
          <MultiLangEditor
            :model="presalePopupTitleI18n"
            :languages="enabledLangs"
            title="预售弹窗标题多语言"
          />
        </el-form-item>
        <el-form-item label="预售弹窗内容(多语言):" prop="presalePopupContent" v-if="formData.isPresale && formData.presalePopupEnabled">
          <el-input v-model="formData.presalePopupContent" type="textarea" :rows="3" placeholder="默认预售弹窗内容" class="mb-2" />
          <MultiLangEditor
            :model="presalePopupContentI18n"
            :languages="enabledLangs"
            title="预售弹窗内容多语言"
            input-type="textarea"
            :rows="3"
          />
        </el-form-item>

        <h4 class="flex justify-between items-center">
          <span>商品规格(可作为商品选项的属性，例：尺码，颜色)</span>
          <div>
            <el-dropdown v-if="specDictList.length" @command="selectSpecFromDict" style="margin-right:8px;">
              <el-button type="success" size="small">从字典选择<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-for="d in specDictList" :key="d.ID" :command="d">{{ d.label }} - {{ d.value }}</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button
              type="primary"
              icon="plus"
              size="small"
              @click="formData.specs.push({
                name: '',
                nameI18n: {},
                value: '',
                valueI18n: {}
              })"
            >手动添加</el-button>
          </div>
        </h4>
        <div
          v-for="(item,index) in formData.specs"
          :key="'spec'+index"
          class="mb-2 p-2" style="border: 1px solid #ebeef5; border-radius: 4px;"
        >
          <el-row :gutter="12">
            <el-col :span="8">
              <el-form-item label="规格名称(默认):" class="mb-0">
                <el-input v-model="item.name" :clearable="true" placeholder="默认规格名称" />
              </el-form-item>
              <div class="mt-2">
                <MultiLangEditor
                  :model="item.nameI18n"
                  :languages="enabledLangs"
                  :title="'规格名称多语言 #' + (index + 1)"
                />
              </div>
            </el-col>
            <el-col :span="8">
              <el-form-item label="规格编码(默认):" class="mb-0">
                <el-input v-model="item.value" :clearable="true" placeholder="请输入规格编码" />
              </el-form-item>
              <div class="mt-2">
                <MultiLangEditor
                  :model="item.valueI18n"
                  :languages="enabledLangs"
                  :title="'规格编码多语言 #' + (index + 1)"
                />
              </div>
            </el-col>
            <el-col :span="8">
              <el-form-item label="操作:">
                <el-button type="danger" icon="delete" @click="formData.specs.splice(index,1)">删除</el-button>
              </el-form-item>
            </el-col>
          </el-row>
        </div>
        <h4 class="flex justify-between items-center">
          <span>商品属性(仅作为展示属性，例：厂商，材料)</span>
          <div>
            <el-dropdown v-if="attrDictList.length" @command="selectAttrFromDict" style="margin-right:8px;">
              <el-button type="success" size="small">从字典选择<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-for="d in attrDictList" :key="d.ID" :command="d">{{ d.label }} - {{ d.value }}</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button
              type="primary"
              icon="plus"
              size="small"
              @click="formData.attrs.push({
                name: '',
                nameI18n: {},
                value: '',
                valueI18n: {}
              })"
            >手动添加</el-button>
          </div>
        </h4>
        <div
          v-for="(item,index) in formData.attrs"
          :key="'attr'+index"
          class="mb-2 p-2" style="border: 1px solid #ebeef5; border-radius: 4px;"
        >
          <el-row :gutter="12">
            <el-col :span="8">
              <el-form-item label="属性名称(默认):" class="mb-0">
                <el-input v-model="item.name" :clearable="true" placeholder="默认属性名称" />
              </el-form-item>
              <div class="mt-2">
                <MultiLangEditor
                  :model="item.nameI18n"
                  :languages="enabledLangs"
                  :title="'属性名称多语言 #' + (index + 1)"
                />
              </div>
            </el-col>
            <el-col :span="8">
              <el-form-item label="属性编码(默认):" class="mb-0">
                <el-input v-model="item.value" :clearable="true" placeholder="请输入属性编码" />
              </el-form-item>
              <div class="mt-2">
                <MultiLangEditor
                  :model="item.valueI18n"
                  :languages="enabledLangs"
                  :title="'属性编码多语言 #' + (index + 1)"
                />
              </div>
            </el-col>
            <el-col :span="8">
              <el-form-item label="操作:">
                <el-button type="danger" icon="delete" @click="formData.attrs.splice(index,1)">删除</el-button>
              </el-form-item>
            </el-col>
          </el-row>
        </div>
        <el-form-item
            label="商品详情:"
            prop="detail"
        >
          <MultiLangEditor
            :model="detailI18n"
            :languages="enabledLangs"
            title="商品详情多语言"
            :use-tabs="true"
            tab-type="card"
          >
            <template #editor="{ lang }">
              <div class="h-[600px]">
                <rich-edit
                  v-model="detailI18n[lang.code]"
                  :clearable="true"
                  :placeholder="`请输入 ${(lang.name || lang.code)} 商品详情`"
                  :upload-folder="GOOD_DETAIL_UPLOAD_FOLDER"
                />
              </div>
            </template>
          </MultiLangEditor>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-dialog
      v-model="priceRateDialogVisible"
      title="商品价格汇率换算"
      width="980px"
      destroy-on-close
    >
      <div class="mb-3 flex flex-wrap items-center gap-2">
        <el-checkbox :model-value="isAllPriceRowsSelected()" @change="toggleSelectAllPriceRows">全选</el-checkbox>
        <el-input-number
          v-model="priceAdjustPercent"
          :precision="2"
          :step="0.5"
          placeholder="增减百分比"
          style="width: 140px"
        />
        <span class="text-xs text-gray-500">增加%（可负数），按人民币价格换算后叠加</span>
      </div>

      <div class="mb-3 flex flex-wrap gap-2">
        <el-button type="primary" :loading="priceRateLoading" @click="refreshSelectedExchangeRates">更新汇率（选中）</el-button>
        <el-button type="success" @click="applySelectedRateToPrices">补齐（更新选中价格）</el-button>
        <el-button @click="fillMissingLangSlots">补齐多语言槽位</el-button>
      </div>

      <div class="mb-2 text-xs text-gray-500">
        汇率来源：{{ exchangeRateSource || '-' }}
        <span v-if="exchangeRateFetchedAt">，更新时间：{{ exchangeRateFetchedAt }}</span>
      </div>

      <el-table :data="priceRateRows" border max-height="420px">
        <el-table-column label="选择" width="70">
          <template #default="scope">
            <el-checkbox v-model="scope.row.selected" />
          </template>
        </el-table-column>
        <el-table-column label="语言" width="180">
          <template #default="scope">
            {{ scope.row.name }} ({{ scope.row.code }})
          </template>
        </el-table-column>
        <el-table-column label="币种" width="140">
          <template #default="scope">
            <el-input v-model="scope.row.currency" maxlength="3" placeholder="USD" @change="normalizeRowCurrency(scope.row)" />
          </template>
        </el-table-column>
        <el-table-column label="汇率(CNY->币种)" width="190">
          <template #default="scope">
            <el-input-number
              v-model="scope.row.rate"
              :precision="6"
              :min="0"
              :step="0.001"
              controls-position="right"
              style="width: 100%"
            />
          </template>
        </el-table-column>
        <el-table-column label="目标价格(分)">
          <template #default="scope">
            <div class="flex items-center gap-2">
              <el-input-number
                v-model="scope.row.targetPriceFen"
                :precision="0"
                :step="1"
                :min="0"
                controls-position="right"
                style="width: 160px"
                @change="syncRowPriceToI18n(scope.row)"
              />
              <span class="text-xs text-gray-500">≈ {{ formatFen(scope.row.targetPriceFen) }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="scope">
            <el-button link type="primary" @click="refreshSingleRate(scope.row)">更新汇率</el-button>
          </template>
        </el-table-column>
      </el-table>

      <template #footer>
        <el-button @click="priceRateDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import {
  createGood,
  deleteGood,
  deleteGoodByIds,
  updateGood,
  findGood,
  getGoodList
} from '@/api/shop/good'

import {
  getTagList
} from '@/api/shop/tag'

import {
  getCategoryList
} from '@/api/shop/category'
import { getAllSkuSpecs } from '@/api/shop/skuSpec'
import { getEnabledLanguages, getLanguageList } from '@/api/client/language'
import { getDefaultDomain } from '@/api/client/externalLinkDomain'
import {
  DEFAULT_LANG_CURRENCY_MAP,
  fetchExchangeRates,
  getExchangeRateSnapshot,
  calcConvertedFenFromCny,
  normalizePriceI18nMap,
  ensurePriceSlots,
  serializePriceI18nMap,
  formatFen
} from '@/utils/exchange-rate'
import { getUrl } from '@/utils/image'
// 图片选择组件
import SelectImage from '@/components/selectImage/selectImage.vue'

// 全量引入格式化工具 请按需保留
import { formatDate, formatBoolean } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'
import RichEdit from "@/components/richtext/rich-edit.vue"
import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'
import { validateRichTextI18nStructure } from '@/utils/richtext-i18n'

defineOptions({
  name: 'Good'
})

const router = useRouter()
const GOOD_UPLOAD_FOLDER = 'cloth-on/web-cloth/goods'
const GOOD_BANNER_UPLOAD_FOLDER = 'cloth-on/web-cloth/goods/lunbo'
const GOOD_DETAIL_UPLOAD_FOLDER = 'cloth-on/web-cloth/goods/detail'

// === 外部链接域名 ===
const extDomain = ref('')
const loadExtDomain = async () => {
  try {
    const res = await getDefaultDomain()
    if (res.code === 0 && res.data) {
      extDomain.value = res.data.replace(/\/+$/, '')
    }
  } catch (e) { /* ignore */ }
}
const resolveExtUrl = (path) => {
  if (!path) return ''
  if (path.startsWith('http')) return path
  if (extDomain.value) {
    const sep = path.startsWith('/') ? '' : '/'
    return extDomain.value + sep + path
  }
  return path
}

const imageSizeTextMap = ref({})
const imageSizeCache = new Map()
const imageSizePromiseCache = new Map()

const getPrimaryImageLink = (row) => {
  if (row?.externalImagePath) return resolveExtUrl(row.externalImagePath)
  return getUrl(row?.imageUrl)
}

const formatFileSize = (bytes) => {
  const value = Number(bytes)
  if (!Number.isFinite(value) || value <= 0) return ''
  if (value < 1024) return `${Math.round(value)} B`
  if (value < 1024 * 1024) return `${(value / 1024).toFixed(2)} KB`
  return `${(value / (1024 * 1024)).toFixed(2)} MB`
}

const loadImageSizeText = (url) => {
  const target = String(url || '').trim()
  if (!target) return Promise.resolve('')

  if (imageSizeCache.has(target)) {
    return Promise.resolve(imageSizeCache.get(target))
  }

  if (imageSizePromiseCache.has(target)) {
    return imageSizePromiseCache.get(target)
  }

  const pending = fetch(target, { method: 'GET' })
    .then((res) => {
      if (!res.ok) return ''
      return res.blob().then((blob) => formatFileSize(blob?.size))
    })
    .catch(() => '')
    .then((text) => {
      const finalText = text || ''
      imageSizeCache.set(target, finalText)
      imageSizePromiseCache.delete(target)
      return finalText
    })

  imageSizePromiseCache.set(target, pending)
  return pending
}

const warmImageSize = (row) => {
  const id = row?.ID
  if (!id) return
  if (imageSizeTextMap.value[id]) return

  const link = getPrimaryImageLink(row)
  if (!link) return

  loadImageSizeText(link).then((text) => {
    if (!text) return
    imageSizeTextMap.value[id] = text
  })
}

const getImageSizeText = (row) => {
  const id = row?.ID
  if (!id) return ''
  if (imageSizeTextMap.value[id]) return imageSizeTextMap.value[id]
  warmImageSize(row)
  return ''
}

// === SKU规格字典 ===
const specDictList = ref([])
const attrDictList = ref([])
const loadSkuSpecDict = async () => {
  try {
    const [specRes, attrRes] = await Promise.all([getAllSkuSpecs('spec'), getAllSkuSpecs('attr')])
    if (specRes.code === 0) specDictList.value = specRes.data || []
    if (attrRes.code === 0) attrDictList.value = attrRes.data || []
  } catch(e) { /* ignore */ }
}
const selectSpecFromDict = (dictItem) => {
  if (!formData.value.specs) formData.value.specs = []
  const nameI18n = {}
  try { Object.assign(nameI18n, typeof dictItem.labelI18n === 'string' ? JSON.parse(dictItem.labelI18n || '{}') : (dictItem.labelI18n || {})) } catch { /* ignore parse error */ }
  const valueI18n = {}
  try { Object.assign(valueI18n, typeof dictItem.valueI18n === 'string' ? JSON.parse(dictItem.valueI18n || '{}') : (dictItem.valueI18n || {})) } catch { /* ignore parse error */ }
  formData.value.specs.push({ name: dictItem.label, nameI18n, value: dictItem.value, valueI18n })
}
const selectAttrFromDict = (dictItem) => {
  if (!formData.value.attrs) formData.value.attrs = []
  const nameI18n = {}
  try { Object.assign(nameI18n, typeof dictItem.labelI18n === 'string' ? JSON.parse(dictItem.labelI18n || '{}') : (dictItem.labelI18n || {})) } catch { /* ignore parse error */ }
  const valueI18n = {}
  try { Object.assign(valueI18n, typeof dictItem.valueI18n === 'string' ? JSON.parse(dictItem.valueI18n || '{}') : (dictItem.valueI18n || {})) } catch { /* ignore parse error */ }
  formData.value.attrs.push({ name: dictItem.label, nameI18n, value: dictItem.value, valueI18n })
}

// === i18n 多语言支持 ===
const enabledLangs = ref([])
const allLangs = ref([])
const titleI18n = ref({})
const descI18n = ref({})
const detailI18n = ref({})
const presalePopupTitleI18n = ref({})
const presalePopupContentI18n = ref({})

const uniqueLangItems = (items = []) => {
  const map = new Map()
  items.forEach((item) => {
    const code = String(item?.code || '').trim()
    if (!code || map.has(code)) return
    map.set(code, {
      ...item,
      code
    })
  })
  return Array.from(map.values())
}

const loadLangs = async () => {
  try {
    const [enabledRes, listRes] = await Promise.all([
      getEnabledLanguages().catch(() => null),
      getLanguageList({ page: 1, pageSize: 500 }).catch(() => null)
    ])

    if (enabledRes?.code === 0) {
      enabledLangs.value = uniqueLangItems(enabledRes.data || [])
    }

    if (listRes?.code === 0) {
      const listPayload = Array.isArray(listRes?.data?.list)
        ? listRes.data.list
        : (Array.isArray(listRes?.data) ? listRes.data : [])
      allLangs.value = uniqueLangItems(listPayload)
    }

    if (!allLangs.value.length) {
      allLangs.value = [...enabledLangs.value]
    }
  } catch (e) { /* ignore */ }
}

const parseI18nJson = (jsonStr) => {
  if (!jsonStr) return {}
  try { return JSON.parse(jsonStr) } catch { return {} }
}

const serializeI18nJson = (obj) => {
  const filtered = {}
  for (const [k, v] of Object.entries(obj)) {
    if (v) filtered[k] = v
  }
  return Object.keys(filtered).length ? JSON.stringify(filtered) : ''
}

const priceI18n = ref({})
const priceRateDialogVisible = ref(false)
const priceRateRows = ref([])
const priceRateLoading = ref(false)
const priceAdjustPercent = ref(0)
const exchangeRateSource = ref('')
const exchangeRateFetchedAt = ref('')

const getEnabledLangCodes = () => {
  const source = allLangs.value.length ? allLangs.value : enabledLangs.value
  return source
    .map((item) => String(item?.code || '').trim())
    .filter(Boolean)
}

const normalizeCurrencyCode = (value) => {
  return String(value || '').trim().toUpperCase().slice(0, 3)
}

const syncPriceI18nToForm = () => {
  formData.value.priceI18n = serializePriceI18nMap(priceI18n.value)
}

const ensureTextLangSlots = (payload, fallbackValue = '') => {
  const next = { ...(payload || {}) }
  getEnabledLangCodes().forEach((code) => {
    if (!Object.prototype.hasOwnProperty.call(next, code)) {
      next[code] = fallbackValue
    }
  })
  return next
}

const ensurePriceLangSlots = () => {
  const fallbackFen = Math.round(Number(formData.value.price) || 0)
  priceI18n.value = ensurePriceSlots(priceI18n.value, getEnabledLangCodes(), fallbackFen)
  syncPriceI18nToForm()
}

const buildPriceRateRows = () => {
  ensurePriceLangSlots()
  const previousMap = {}
  priceRateRows.value.forEach((row) => {
    previousMap[row.code] = row
  })

  const fallbackFen = Math.round(Number(formData.value.price) || 0)
  const source = allLangs.value.length ? allLangs.value : enabledLangs.value
  priceRateRows.value = source.map((lang) => {
    const code = String(lang?.code || '').trim()
    const previous = previousMap[code] || {}
    const currency = normalizeCurrencyCode(previous.currency || DEFAULT_LANG_CURRENCY_MAP[code] || 'USD')
    const targetPriceFen = Number.isFinite(Number(priceI18n.value[code]))
      ? Number(priceI18n.value[code])
      : fallbackFen
    const defaultRate = code === 'zh' ? 1 : 0

    return {
      code,
      name: lang?.name || code,
      currency,
      rate: Number.isFinite(Number(previous.rate)) ? Number(previous.rate) : defaultRate,
      selected: previous.selected !== false,
      targetPriceFen
    }
  })

  const snapshot = getExchangeRateSnapshot({ base: 'CNY' })
  if (!snapshot?.rates) {
    return
  }

  exchangeRateSource.value = snapshot.source || ''
  exchangeRateFetchedAt.value = snapshot.fetchedAt || ''
  priceRateRows.value.forEach((row) => {
    if (row.code === 'zh') {
      row.rate = 1
      return
    }
    const currency = normalizeCurrencyCode(row.currency)
    const snapshotRate = Number(snapshot.rates?.[currency])
    if (Number.isFinite(snapshotRate) && snapshotRate > 0) {
      row.rate = snapshotRate
    }
  })
}

const openPriceRateDialog = async() => {
  if (!getEnabledLangCodes().length) {
    ElMessage.warning('请先在语言管理中配置语言')
    return
  }
  buildPriceRateRows()
  priceRateDialogVisible.value = true
  await refreshRatesForRows(priceRateRows.value, {
    silent: true,
    forceRefresh: false
  })
}

const isAllPriceRowsSelected = () => {
  return priceRateRows.value.length > 0 && priceRateRows.value.every((row) => row.selected)
}

const toggleSelectAllPriceRows = (checked) => {
  priceRateRows.value.forEach((row) => {
    row.selected = !!checked
  })
}

const normalizeRowCurrency = (row) => {
  row.currency = normalizeCurrencyCode(row.currency)
}

const syncRowPriceToI18n = (row) => {
  const priceFen = Math.max(0, Math.round(Number(row.targetPriceFen) || 0))
  row.targetPriceFen = priceFen
  priceI18n.value[row.code] = priceFen
  syncPriceI18nToForm()
}

const refreshRatesForRows = async(rows, {
  silent = false,
  forceRefresh = false
} = {}) => {
  const validRows = rows.filter((row) => row?.code)
  if (!validRows.length) {
    if (!silent) {
      ElMessage.warning('请先选择需要更新汇率的语言')
    }
    return
  }

  const currencies = validRows
    .map((row) => normalizeCurrencyCode(row.currency))
    .filter(Boolean)

  try {
    priceRateLoading.value = true
    const result = await fetchExchangeRates({
      base: 'CNY',
      currencies,
      forceRefresh
    })
    exchangeRateSource.value = result.source || ''
    exchangeRateFetchedAt.value = result.fetchedAt || ''

    validRows.forEach((row) => {
      if (row.code === 'zh') {
        row.rate = 1
        return
      }
      const currency = normalizeCurrencyCode(row.currency)
      const nextRate = Number(result?.rates?.[currency])
      if (Number.isFinite(nextRate) && nextRate > 0) {
        row.rate = nextRate
      }
    })
    if (!silent) {
      ElMessage.success('汇率更新完成')
    }
  } catch (error) {
    if (!silent) {
      ElMessage.error(error?.message || '汇率更新失败')
    }
  } finally {
    priceRateLoading.value = false
  }
}

const refreshSelectedExchangeRates = async() => {
  const selectedRows = priceRateRows.value.filter((row) => row.selected)
  await refreshRatesForRows(selectedRows, { forceRefresh: true })
}

const refreshSingleRate = async(row) => {
  await refreshRatesForRows([row], { forceRefresh: true })
}

const applySelectedRateToPrices = () => {
  const selectedRows = priceRateRows.value.filter((row) => row.selected)
  if (!selectedRows.length) {
    ElMessage.warning('请至少选择一个语言')
    return
  }

  let updatedCount = 0
  selectedRows.forEach((row) => {
    const convertedFen = calcConvertedFenFromCny(formData.value.price, row.rate, priceAdjustPercent.value)
    if (convertedFen === null) {
      return
    }
    row.targetPriceFen = convertedFen
    priceI18n.value[row.code] = convertedFen
    updatedCount += 1
  })

  syncPriceI18nToForm()
  ElMessage.success(`已更新 ${updatedCount} 个语言价格`)
}

const fillMissingLangSlots = () => {
  if (!getEnabledLangCodes().length) {
    ElMessage.warning('请先在语言管理中配置语言')
    return
  }

  titleI18n.value = ensureTextLangSlots(titleI18n.value, formData.value.title || '')
  descI18n.value = ensureTextLangSlots(descI18n.value, formData.value.description || '')
  detailI18n.value = ensureTextLangSlots(detailI18n.value, formData.value.detail || '')
  presalePopupTitleI18n.value = ensureTextLangSlots(presalePopupTitleI18n.value, formData.value.presalePopupTitle || '')
  presalePopupContentI18n.value = ensureTextLangSlots(presalePopupContentI18n.value, formData.value.presalePopupContent || '')

  formData.value.banner = (formData.value.banner || []).map((item) => {
    const nextTextI18n = ensureTextLangSlots(item.textI18n || parseI18nJson(item.text), item.text || '')
    return {
      ...item,
      textI18n: nextTextI18n
    }
  })

  formData.value.specs = (formData.value.specs || []).map((item) => {
    return {
      ...item,
      nameI18n: ensureTextLangSlots(item.nameI18n || parseI18nJson(item.name), item.name || ''),
      valueI18n: ensureTextLangSlots(item.valueI18n || parseI18nJson(item.value), item.value || '')
    }
  })

  formData.value.attrs = (formData.value.attrs || []).map((item) => {
    return {
      ...item,
      nameI18n: ensureTextLangSlots(item.nameI18n || parseI18nJson(item.name), item.name || ''),
      valueI18n: ensureTextLangSlots(item.valueI18n || parseI18nJson(item.value), item.value || '')
    }
  })

  ensurePriceLangSlots()
  buildPriceRateRows()
  ElMessage.success('多语言槽位补齐完成')
}

onMounted(() => { loadLangs(); loadExtDomain(); loadSkuSpecDict() })

const tableLang = ref('')

const formatI18nField = (val) => {
  if (!val) return ''
  if (typeof val === 'string' && val.startsWith('{')) {
    try {
      const obj = JSON.parse(val)
      if (tableLang.value && obj[tableLang.value]) return obj[tableLang.value]
      return obj['zh'] || Object.values(obj)[0] || val
    } catch { return val }
  }
  return val
}

const getLocalizedPriceFen = (basePrice, i18nPayload, langCode) => {
  const map = normalizePriceI18nMap(i18nPayload)
  if (langCode && Number.isFinite(Number(map[langCode]))) {
    return Number(map[langCode])
  }
  const base = Number(basePrice)
  return Number.isFinite(base) ? base : 0
}

const formatTablePrice = (row) => {
  return formatFen(getLocalizedPriceFen(row?.price, row?.priceI18n, tableLang.value))
}

const setSKU = (row) => {
  router.push({ name: 'sku', query: { id: row.ID }})
}

const gotoPurchase = (row) => {
  router.push({ name: 'goodPurchase', query: { goodID: row.ID }})
}

const getTotalInventory = (row) => {
  if (!row.skus || !row.skus.length) return 0
  return row.skus.reduce((sum, sku) => sum + (sku.inventory || 0), 0)
}

const tags = ref([])

const getTags = async() => {
  const res = await getTagList()
  if (res.code === 0) {
    tags.value = res.data.list || []
  }
}

getTags()

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  description: '',
  imageUrl: '',
  upperImage: '',
  lowerImage: '',
  externalImagePath: '',
  banner: [],
  price: 0,
  priceI18n: '',
  rating: 0,
  reviewCount: 0,
  title: '',
  categoryID: undefined,
  status: false,
  recommend: false,
  postage: 0,
  discount: 0,
  specs: [],
  attrs: [],
  detail: '',
  remark: '',
  pointsEnabled: false,
  pointsMaxUse: 0,
  pointsUseTimes: 1,
  isPresale: false,
  presaleQty: 0,
  presaleStart: null,
  presaleEnd: null,
  presaleEnabled: true,
  presaleSort: 0,
  presalePopupEnabled: false,
  presalePopupTitle: '',
  presalePopupContent: '',
})
const categoryList = ref([])
const getCategoryListFunc = async() => {
  const res = await getCategoryList()
  if (res.code === 0) {
    categoryList.value = res.data.list || []
  }
}

getCategoryListFunc()

// 验证规则
const rule = reactive({
})

const searchRule = reactive({
  createdAt: [
    { validator: (rule, value, callback) => {
      if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
        callback(new Error('请填写结束日期'))
      } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
        callback(new Error('请填写开始日期'))
      } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
        callback(new Error('开始日期应当早于结束日期'))
      } else {
        callback()
      }
    }, trigger: 'change' }
  ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    if (searchInfo.value.status === '') {
      searchInfo.value.status = null
    }
    if (searchInfo.value.recommend === '') {
      searchInfo.value.recommend = null
    }
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  const table = await getGoodList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    imageSizeTextMap.value = {}
    tableData.value.forEach((row) => warmImageSize(row))
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async() => {
}

// 获取需要的字典 可能为空 按需保留
setOptions()

// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    deleteGoodFunc(row)
  })
}

// 多选删除
const onDelete = async() => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    const IDs = []
    if (multipleSelection.value.length === 0) {
      ElMessage({
        type: 'warning',
        message: '请选择要删除的数据'
      })
      return
    }
    multipleSelection.value &&
        multipleSelection.value.map(item => {
          IDs.push(item.ID)
        })
    const res = await deleteGoodByIds({ IDs })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
      if (tableData.value.length === IDs.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateGoodFunc = async(row) => {
  const res = await findGood({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.regood
    if (!formData.value.upperImage) {
      formData.value.upperImage = ''
    }
    if (!formData.value.lowerImage) {
      formData.value.lowerImage = ''
    }
    if (!formData.value.specs) {
      formData.value.specs = []
    }
    if (!formData.value.attrs) {
      formData.value.attrs = []
    }
    if (!formData.value.banner) {
      formData.value.banner = []
    }
    // 兼容旧数据：如果 banner 是简单字符串数组，转为新结构
    if (formData.value.banner.length > 0 && typeof formData.value.banner[0] === 'string') {
      formData.value.banner = formData.value.banner.map(url => ({
        url, type: 'image', externalUrl: '', text: '', textI18n: {}, textColor: '#FFFFFF', textSize: 14, textPosition: 'bottom'
      }))
    }
    // 解析 banner 文字 i18n
    formData.value.banner.forEach(item => {
      if (!item.textI18n) {
        item.textI18n = parseI18nJson(item.text)
      }
    })
    // 解析 specs/attrs 名称和编码 i18n
    formData.value.specs.forEach(item => {
      item.nameI18n = parseI18nJson(item.name)
      if (!item.valueI18n) item.valueI18n = parseI18nJson(item.value)
      else if (typeof item.valueI18n === 'string') item.valueI18n = parseI18nJson(item.valueI18n)
    })
    formData.value.attrs.forEach(item => {
      item.nameI18n = parseI18nJson(item.name)
      if (!item.valueI18n) item.valueI18n = parseI18nJson(item.value)
      else if (typeof item.valueI18n === 'string') item.valueI18n = parseI18nJson(item.valueI18n)
    })
    // 解析 i18n 字段
    titleI18n.value = parseI18nJson(formData.value.title)
    descI18n.value = parseI18nJson(formData.value.description)
    detailI18n.value = parseI18nJson(formData.value.detail)
    presalePopupTitleI18n.value = parseI18nJson(formData.value.presalePopupTitle)
    presalePopupContentI18n.value = parseI18nJson(formData.value.presalePopupContent)
    priceI18n.value = normalizePriceI18nMap(formData.value.priceI18n)
    ensurePriceLangSlots()
    buildPriceRateRows()
    dialogFormVisible.value = true
  }
}

// 删除行
const deleteGoodFunc = async(row) => {
  const res = await deleteGood({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除成功'
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  titleI18n.value = {}
  descI18n.value = {}
  detailI18n.value = {}
  presalePopupTitleI18n.value = {}
  presalePopupContentI18n.value = {}
  priceI18n.value = {}
  priceRateRows.value = []
  priceAdjustPercent.value = 0
  exchangeRateSource.value = ''
  exchangeRateFetchedAt.value = ''
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  priceRateDialogVisible.value = false
  titleI18n.value = {}
  descI18n.value = {}
  detailI18n.value = {}
  presalePopupTitleI18n.value = {}
  presalePopupContentI18n.value = {}
  priceI18n.value = {}
  priceRateRows.value = []
  priceAdjustPercent.value = 0
  exchangeRateSource.value = ''
  exchangeRateFetchedAt.value = ''
  formData.value = {
    description: '',
    imageUrl: '',
    upperImage: '',
    lowerImage: '',
    externalImagePath: '',
    banner: [],
    price: 0,
    priceI18n: '',
    rating: 0,
    reviewCount: 0,
    title: '',
    categoryID: undefined,
    status: false,
    recommend: false,
    postage: 0,
    specs: [],
    attrs: [],
    discount: 0,
    detail: '',
    remark: '',
    pointsEnabled: false,
    pointsMaxUse: 0,
    pointsUseTimes: 1,
    isPresale: false,
    presaleQty: 0,
    presaleStart: null,
    presaleEnd: null,
    presaleEnabled: true,
    presaleSort: 0,
    presalePopupEnabled: false,
    presalePopupTitle: '',
    presalePopupContent: '',
  }
}
// 弹窗确定
const enterDialog = async() => {
  elFormRef.value?.validate(async(valid) => {
    if (!valid) return

    // 统一富文本多语校验钩子：翻译后保持关键 HTML 结构一致，避免链接/列表/媒体标签被破坏。
    const richTextValidation = validateRichTextI18nStructure(detailI18n.value, {
      fieldLabel: '商品详情',
      preferredBaseLang: 'zh',
    })
    if (!richTextValidation.valid) {
      ElMessage.warning(richTextValidation.message)
      return
    }

    // 构建深拷贝 payload，避免序列化时污染 formData（导致保存失败后 nameI18n 丢失）
    const payload = JSON.parse(JSON.stringify(formData.value))
    payload.priceI18n = serializePriceI18nMap(priceI18n.value)
    // 序列化 i18n 字段（写 payload，不写 formData）
    if (getEnabledLangCodes().length) {
      payload.title = serializeI18nJson(titleI18n.value) || payload.title
      payload.description = serializeI18nJson(descI18n.value) || payload.description
      payload.detail = serializeI18nJson(detailI18n.value) || payload.detail
      payload.presalePopupTitle = serializeI18nJson(presalePopupTitleI18n.value) || payload.presalePopupTitle
      payload.presalePopupContent = serializeI18nJson(presalePopupContentI18n.value) || payload.presalePopupContent
      payload.banner.forEach(item => {
        if (item.textI18n) {
          item.text = serializeI18nJson(item.textI18n) || item.text
          delete item.textI18n
        }
      })
      payload.specs.forEach(item => {
        if (item.nameI18n) {
          item.name = serializeI18nJson(item.nameI18n) || item.name
          delete item.nameI18n
        }
        if (item.valueI18n) {
          item.value = serializeI18nJson(item.valueI18n) || item.value
          delete item.valueI18n
        }
      })
      payload.attrs.forEach(item => {
        if (item.nameI18n) {
          item.name = serializeI18nJson(item.nameI18n) || item.name
          delete item.nameI18n
        }
        if (item.valueI18n) {
          item.value = serializeI18nJson(item.valueI18n) || item.value
          delete item.valueI18n
        }
      })
    }
    let res
    switch (type.value) {
      case 'create':
        res = await createGood(payload)
        break
      case 'update':
        res = await updateGood(payload)
        break
      default:
        res = await createGood(payload)
        break
    }
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '创建/更改成功'
      })
      closeDialog()
      getTableData()
    }
  })
}

const onStatusChange = async(row) => {
  const res = await updateGood(row)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '状态更新成功'
    })
  } else {
    row.status = !row.status
  }
}
</script>

<style>
.good-image-wrap {
  position: relative;
  width: 100px;
  height: 100px;
}

.good-image {
  width: 100px;
  height: 100px;
}

.good-image-size-badge {
  position: absolute;
  right: 4px;
  top: 4px;
  padding: 2px 6px;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.7);
  color: #fff;
  font-size: 11px;
  line-height: 1.2;
}

</style>
