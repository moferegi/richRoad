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
            <el-image
              style="width: 100px; height: 100px"
              :src="scope.row.externalImagePath ? resolveExtUrl(scope.row.externalImagePath) : getUrl(scope.row.imageUrl)"
              fit="cover"
            />
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="商品价格"
          prop="price"
          width="120"
        />
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
            <div v-if="enabledLangs.length" class="pl-2">
              <div v-for="lang in enabledLangs" :key="lang.code" class="flex items-center mb-1">
                <span class="w-14 text-xs text-right mr-1">{{ lang.code }}:</span>
                <el-input v-model="titleI18n[lang.code]" :placeholder="`${lang.name}`" size="small" />
              </div>
            </div>
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
                    <SelectImage v-model="item.url" :file-type="item.type === 'video' ? 'video' : 'image'" />
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
                    <div v-if="enabledLangs.length" class="w-full mt-1">
                      <div v-for="lang in enabledLangs" :key="lang.code" class="flex items-center mb-1">
                        <span class="w-14 text-xs text-right mr-1">{{ lang.code }}:</span>
                        <el-input v-model="item.textI18n[lang.code]" :placeholder="lang.name" size="small" />
                      </div>
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
          <div v-if="enabledLangs.length" class="w-full mt-2">
            <div v-for="lang in enabledLangs" :key="lang.code" class="flex items-center mb-1">
              <span class="w-14 text-xs text-right mr-1">{{ lang.code }}:</span>
              <el-input v-model="descI18n[lang.code]" :placeholder="`${lang.name}`" size="small" />
            </div>
          </div>
        </el-form-item>

        <el-row :gutter="12">
          <el-col :span="8">
            <el-form-item
              label="商品价格(单位：分):"
              prop="price"
            >
              <el-input-number
                v-model="formData.price"
                style="width:100%"
                :precision="2"
                :clearable="true"
              />
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
          <el-tabs type="border-card" style="width:100%;">
            <el-tab-pane v-for="lang in enabledLangs" :key="lang.code" :label="lang.name">
              <el-input v-model="presalePopupTitleI18n[lang.code]" :placeholder="`${lang.name} 弹窗标题`" />
            </el-tab-pane>
          </el-tabs>
        </el-form-item>
        <el-form-item label="预售弹窗内容(多语言):" prop="presalePopupContent" v-if="formData.isPresale && formData.presalePopupEnabled">
          <el-tabs type="border-card" style="width:100%;">
            <el-tab-pane v-for="lang in enabledLangs" :key="lang.code" :label="lang.name">
              <el-input v-model="presalePopupContentI18n[lang.code]" type="textarea" :rows="3" :placeholder="`${lang.name} 弹窗内容`" />
            </el-tab-pane>
          </el-tabs>
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
              <div v-if="enabledLangs.length && item.nameI18n" class="pl-2 mt-1">
                <div v-for="lang in enabledLangs" :key="lang.code" class="flex items-center mb-1">
                  <span class="w-14 text-xs text-right mr-1">{{ lang.code }}:</span>
                  <el-input v-model="item.nameI18n[lang.code]" :placeholder="lang.name" size="small" />
                </div>
              </div>
            </el-col>
            <el-col :span="8">
              <el-form-item label="规格编码(默认):" class="mb-0">
                <el-input v-model="item.value" :clearable="true" placeholder="请输入规格编码" />
              </el-form-item>
              <div v-if="enabledLangs.length && item.valueI18n" class="pl-2 mt-1">
                <div v-for="lang in enabledLangs" :key="lang.code" class="flex items-center mb-1">
                  <span class="w-14 text-xs text-right mr-1">{{ lang.code }}:</span>
                  <el-input v-model="item.valueI18n[lang.code]" :placeholder="lang.name" size="small" />
                </div>
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
              <div v-if="enabledLangs.length && item.nameI18n" class="pl-2 mt-1">
                <div v-for="lang in enabledLangs" :key="lang.code" class="flex items-center mb-1">
                  <span class="w-14 text-xs text-right mr-1">{{ lang.code }}:</span>
                  <el-input v-model="item.nameI18n[lang.code]" :placeholder="lang.name" size="small" />
                </div>
              </div>
            </el-col>
            <el-col :span="8">
              <el-form-item label="属性编码(默认):" class="mb-0">
                <el-input v-model="item.value" :clearable="true" placeholder="请输入属性编码" />
              </el-form-item>
              <div v-if="enabledLangs.length && item.valueI18n" class="pl-2 mt-1">
                <div v-for="lang in enabledLangs" :key="lang.code" class="flex items-center mb-1">
                  <span class="w-14 text-xs text-right mr-1">{{ lang.code }}:</span>
                  <el-input v-model="item.valueI18n[lang.code]" :placeholder="lang.name" size="small" />
                </div>
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
          <div v-if="enabledLangs.length" class="w-full">
            <el-tabs v-model="detailLangTab" type="card">
              <el-tab-pane v-for="lang in enabledLangs" :key="lang.code" :label="lang.name || lang.code" :name="lang.code">
                <div class="h-[600px]">
                  <rich-edit
                    v-model="detailI18n[lang.code]"
                    :clearable="true"
                    :placeholder="`请输入 ${lang.name} 商品详情`"
                  />
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
          <div v-else class="h-[800px]">
            <rich-edit
                v-model="formData.detail"
                :clearable="true"
                placeholder="请输入商品描述"
            />
          </div>
        </el-form-item>
      </el-form>
    </el-drawer>
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
import { getEnabledLanguages } from '@/api/client/language'
import { getDefaultDomain } from '@/api/client/externalLinkDomain'
import { getUrl } from '@/utils/image'
// 图片选择组件
import SelectImage from '@/components/selectImage/selectImage.vue'

// 全量引入格式化工具 请按需保留
import { formatDate, formatBoolean } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'
import RichEdit from "@/components/richtext/rich-edit.vue"

defineOptions({
  name: 'Good'
})

const router = useRouter()

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
  try { Object.assign(nameI18n, typeof dictItem.labelI18n === 'string' ? JSON.parse(dictItem.labelI18n || '{}') : (dictItem.labelI18n || {})) } catch {}
  const valueI18n = {}
  try { Object.assign(valueI18n, typeof dictItem.valueI18n === 'string' ? JSON.parse(dictItem.valueI18n || '{}') : (dictItem.valueI18n || {})) } catch {}
  formData.value.specs.push({ name: dictItem.label, nameI18n, value: dictItem.value, valueI18n })
}
const selectAttrFromDict = (dictItem) => {
  if (!formData.value.attrs) formData.value.attrs = []
  const nameI18n = {}
  try { Object.assign(nameI18n, typeof dictItem.labelI18n === 'string' ? JSON.parse(dictItem.labelI18n || '{}') : (dictItem.labelI18n || {})) } catch {}
  const valueI18n = {}
  try { Object.assign(valueI18n, typeof dictItem.valueI18n === 'string' ? JSON.parse(dictItem.valueI18n || '{}') : (dictItem.valueI18n || {})) } catch {}
  formData.value.attrs.push({ name: dictItem.label, nameI18n, value: dictItem.value, valueI18n })
}

// === i18n 多语言支持 ===
const enabledLangs = ref([])
const titleI18n = ref({})
const descI18n = ref({})
const detailI18n = ref({})
const presalePopupTitleI18n = ref({})
const presalePopupContentI18n = ref({})
// 当前详情编辑语言tab
const detailLangTab = ref('')

const loadLangs = async () => {
  try {
    const res = await getEnabledLanguages()
    if (res.code === 0) {
      enabledLangs.value = res.data || []
      if (enabledLangs.value.length > 0) {
        detailLangTab.value = enabledLangs.value[0].code
      }
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
  externalImagePath: '',
  banner: [],
  price: 0,
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
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  titleI18n.value = {}
  descI18n.value = {}
  detailI18n.value = {}
  presalePopupTitleI18n.value = {}
  presalePopupContentI18n.value = {}
  formData.value = {
    description: '',
    imageUrl: '',
    externalImagePath: '',
    banner: [],
    price: 0,
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
    // 构建深拷贝 payload，避免序列化时污染 formData（导致保存失败后 nameI18n 丢失）
    const payload = JSON.parse(JSON.stringify(formData.value))
    // 序列化 i18n 字段（写 payload，不写 formData）
    if (enabledLangs.value.length) {
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
          delete item.valueI18n
        }
      })
      payload.attrs.forEach(item => {
        if (item.nameI18n) {
          item.name = serializeI18nJson(item.nameI18n) || item.name
          delete item.nameI18n
        }
        if (item.valueI18n) {
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

</style>
