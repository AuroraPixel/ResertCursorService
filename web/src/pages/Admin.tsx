import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { toast, Toaster } from 'react-hot-toast';
import axios from 'axios';
import { useAuth } from '../hooks/useAuth';
import { Dialog } from '@headlessui/react';
import { 
  ArrowLeftOnRectangleIcon, 
  ClipboardDocumentIcon,
  CheckCircleIcon,
  XCircleIcon,
  EyeIcon,
  PlusIcon,
  ChevronLeftIcon,
  ChevronRightIcon
} from '@heroicons/react/24/outline';

interface PageItem {
  type: 'page' | 'dots';
  number?: number;
  position?: 'start' | 'middle' | 'end';
}

interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

interface ActivationCode {
  id: number;
  code: string;
  expiresAt: string;
  maxAccounts: number;
  status: 'enabled' | 'disabled';
  createdAt: string;
  accounts: Array<{
    id: number;
    email: string;
    emailPassword: string;
    cursorPassword: string;
    accessToken: string;
    refreshToken: string;
    activationCodeId: number;
  }>;
}

export default function Admin() {
  const navigate = useNavigate();
  const { logout } = useAuth();
  const [codes, setCodes] = useState<ActivationCode[]>([]);
  const [selectedCode, setSelectedCode] = useState<ActivationCode | null>(null);
  const [showModal, setShowModal] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isGenerating, setIsGenerating] = useState(false);
  const [pagination, setPagination] = useState({
    page: 1,
    pageSize: 10,
    total: 0,
    totalPages: 0
  });
  const [accountsPagination, setAccountsPagination] = useState({
    page: 1,
    pageSize: 5
  });
  const [newCode, setNewCode] = useState({
    duration: 5,
    maxAccounts: 1,
  });

  useEffect(() => {
    loadActivationCodes();
  }, [pagination.page]);

  const loadActivationCodes = async () => {
    try {
      setIsLoading(true);
      const response = await axios.get<PaginatedResponse<ActivationCode>>('/api/activation-codes', {
        params: {
          page: pagination.page,
          pageSize: pagination.pageSize
        }
      });
      setCodes(response.data.items);
      setPagination(prev => ({
        ...prev,
        total: response.data.total,
        totalPages: response.data.totalPages
      }));
    } catch (error) {
      console.error('Failed to load activation codes:', error);
      toast.error('加载激活码列表失败');
    } finally {
      setIsLoading(false);
    }
  };

  const generateActivationCode = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setIsGenerating(true);
      const response = await axios.post('/api/activation-codes', {
        duration: newCode.duration,
        maxAccounts: newCode.maxAccounts,
      });
      
      if (response.data.code) {
        toast.success('激活码生成成功');
        await loadActivationCodes();
        setNewCode({ ...newCode, maxAccounts: 1 });
      }
    } catch (error) {
      console.error('Failed to generate activation code:', error);
      toast.error('生成激活码失败，请重试');
    } finally {
      setIsGenerating(false);
    }
  };

  const toggleStatus = async (code: ActivationCode) => {
    try {
      setIsLoading(true);
      await axios.put(`/api/activation-codes/${code.id}/status`, {
        status: code.status === 'enabled' ? 'disabled' : 'enabled',
      });
      toast.success('状态更新成功');
      await loadActivationCodes();
    } catch (error) {
      console.error('Failed to update activation code status:', error);
      toast.error('更新状态失败');
    } finally {
      setIsLoading(false);
    }
  };

  const viewDetails = (code: ActivationCode) => {
    setSelectedCode(code);
    setShowModal(true);
  };

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    toast.success('已复制到剪贴板');
  };

  const handlePageChange = (newPage: number) => {
    setPagination(prev => ({ ...prev, page: newPage }));
  };

  const getPageNumbers = (): PageItem[] => {
    const pages: PageItem[] = [];
    const { page, totalPages } = pagination;
    
    if (totalPages <= 7) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push({ type: 'page', number: i });
      }
    } else {
      if (page <= 3) {
        for (let i = 1; i <= 5; i++) {
          pages.push({ type: 'page', number: i });
        }
        pages.push({ type: 'dots', position: 'middle' });
        pages.push({ type: 'page', number: totalPages });
      } else if (page >= totalPages - 2) {
        pages.push({ type: 'page', number: 1 });
        pages.push({ type: 'dots', position: 'start' });
        for (let i = totalPages - 4; i <= totalPages; i++) {
          pages.push({ type: 'page', number: i });
        }
      } else {
        pages.push({ type: 'page', number: 1 });
        pages.push({ type: 'dots', position: 'start' });
        for (let i = page - 1; i <= page + 1; i++) {
          pages.push({ type: 'page', number: i });
        }
        pages.push({ type: 'dots', position: 'end' });
        pages.push({ type: 'page', number: totalPages });
      }
    }
    
    return pages;
  };

  const renderPagination = () => (
    <div className="flex items-center justify-between px-4 py-3 bg-white border-t border-gray-200 sm:px-6">
      <div className="flex justify-between flex-1 sm:hidden">
        <button
          onClick={() => handlePageChange(pagination.page - 1)}
          disabled={pagination.page === 1}
          className="relative inline-flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          上一页
        </button>
        <button
          onClick={() => handlePageChange(pagination.page + 1)}
          disabled={pagination.page === pagination.totalPages}
          className="relative inline-flex items-center px-4 py-2 ml-3 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          下一页
        </button>
      </div>
      <div className="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
        <div>
          <p className="text-sm text-gray-700">
            显示第 <span className="font-medium">{(pagination.page - 1) * pagination.pageSize + 1}</span> 到{' '}
            <span className="font-medium">
              {Math.min(pagination.page * pagination.pageSize, pagination.total)}
            </span>{' '}
            条，共 <span className="font-medium">{pagination.total}</span> 条
          </p>
        </div>
        <div>
          <nav className="inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
            <button
              onClick={() => handlePageChange(pagination.page - 1)}
              disabled={pagination.page === 1}
              className="relative inline-flex items-center px-2 py-2 text-gray-400 rounded-l-md border border-gray-300 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span className="sr-only">上一页</span>
              <ChevronLeftIcon className="h-5 w-5" aria-hidden="true" />
            </button>
            {getPageNumbers().map((item) => (
              <button
                key={item.type === 'dots' ? `dots-${item.position}` : `page-${item.number}`}
                onClick={() => item.type === 'page' && item.number && handlePageChange(item.number)}
                disabled={item.type === 'dots'}
                className={`relative inline-flex items-center px-4 py-2 text-sm font-medium ${
                  item.type === 'page' && item.number === pagination.page
                    ? 'z-10 bg-indigo-50 border-indigo-500 text-indigo-600'
                    : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
                } border`}
              >
                {item.type === 'dots' ? '...' : item.number}
              </button>
            ))}
            <button
              onClick={() => handlePageChange(pagination.page + 1)}
              disabled={pagination.page === pagination.totalPages}
              className="relative inline-flex items-center px-2 py-2 text-gray-400 rounded-r-md border border-gray-300 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span className="sr-only">下一页</span>
              <ChevronRightIcon className="h-5 w-5" aria-hidden="true" />
            </button>
          </nav>
        </div>
      </div>
    </div>
  );

  const renderAccountsList = () => {
    if (!selectedCode) return null;
    
    const accounts = selectedCode.accounts || [];
    const startIndex = (accountsPagination.page - 1) * accountsPagination.pageSize;
    const endIndex = startIndex + accountsPagination.pageSize;
    const currentAccounts = accounts.slice(startIndex, endIndex);
    const totalPages = Math.ceil(accounts.length / accountsPagination.pageSize);

    return (
      <div className="mt-4">
        <h4 className="text-sm font-medium text-gray-900">账户列表</h4>
        <div className="mt-2 space-y-4">
          {currentAccounts.map((account) => (
            <motion.div
              key={account.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="bg-gray-50 rounded-lg p-4"
            >
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-xs text-gray-500">邮箱</p>
                  <div className="flex items-center mt-1">
                    <p className="text-sm font-medium">{account.email}</p>
                    <button
                      onClick={() => copyToClipboard(account.email)}
                      className="ml-2 text-gray-400 hover:text-gray-600"
                    >
                      <ClipboardDocumentIcon className="h-4 w-4" />
                    </button>
                  </div>
                </div>
                <div>
                  <p className="text-xs text-gray-500">邮箱密码</p>
                  <div className="flex items-center mt-1">
                    <p className="text-sm font-medium">{account.emailPassword}</p>
                    <button
                      onClick={() => copyToClipboard(account.emailPassword)}
                      className="ml-2 text-gray-400 hover:text-gray-600"
                    >
                      <ClipboardDocumentIcon className="h-4 w-4" />
                    </button>
                  </div>
                </div>
              </div>
              <div className="mt-2">
                <p className="text-xs text-gray-500">Cursor密码</p>
                <div className="flex items-center mt-1">
                  <p className="text-sm font-medium break-all">{account.cursorPassword}</p>
                  <button
                    onClick={() => copyToClipboard(account.cursorPassword)}
                    className="ml-2 text-gray-400 hover:text-gray-600 flex-shrink-0"
                  >
                    <ClipboardDocumentIcon className="h-4 w-4" />
                  </button>
                </div>
              </div>
              <div className="mt-2">
                <p className="text-xs text-gray-500">Access Token</p>
                <div className="flex items-center mt-1">
                  <p className="text-sm font-medium break-all">{account.accessToken}</p>
                  <button
                    onClick={() => copyToClipboard(account.accessToken)}
                    className="ml-2 text-gray-400 hover:text-gray-600 flex-shrink-0"
                  >
                    <ClipboardDocumentIcon className="h-4 w-4" />
                  </button>
                </div>
              </div>
              <div className="mt-2">
                <p className="text-xs text-gray-500">Refresh Token</p>
                <div className="flex items-center mt-1">
                  <p className="text-sm font-medium break-all">{account.refreshToken}</p>
                  <button
                    onClick={() => copyToClipboard(account.refreshToken)}
                    className="ml-2 text-gray-400 hover:text-gray-600 flex-shrink-0"
                  >
                    <ClipboardDocumentIcon className="h-4 w-4" />
                  </button>
                </div>
              </div>
            </motion.div>
          ))}
        </div>
        {accounts.length > accountsPagination.pageSize && (
          <div className="mt-4 flex justify-center">
            <nav className="inline-flex -space-x-px rounded-md shadow-sm">
              <button
                onClick={() => setAccountsPagination(prev => ({ ...prev, page: prev.page - 1 }))}
                disabled={accountsPagination.page === 1}
                className="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <ChevronLeftIcon className="h-5 w-5" />
              </button>
              {Array.from({ length: totalPages }).map((_, i) => (
                <button
                  key={`account-page-${i + 1}`}
                  onClick={() => setAccountsPagination(prev => ({ ...prev, page: i + 1 }))}
                  className={`relative inline-flex items-center px-4 py-2 border text-sm font-medium ${
                    i + 1 === accountsPagination.page
                      ? 'z-10 bg-indigo-50 border-indigo-500 text-indigo-600'
                      : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
                  }`}
                >
                  {i + 1}
                </button>
              ))}
              <button
                onClick={() => setAccountsPagination(prev => ({ ...prev, page: prev.page + 1 }))}
                disabled={accountsPagination.page === totalPages}
                className="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <ChevronRightIcon className="h-5 w-5" />
              </button>
            </nav>
          </div>
        )}
      </div>
    );
  };

  const renderTable = () => (
    <div className="bg-white shadow-sm rounded-lg overflow-hidden">
      <div className="px-6 py-4 border-b border-gray-200">
        <h2 className="text-lg font-medium text-gray-900">激活码列表</h2>
      </div>
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                激活码
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                状态
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                创建时间
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                过期时间
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                账户使用情况
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                操作
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {codes?.map((code) => (
              <motion.tr
                key={code.id}
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                whileHover={{ backgroundColor: 'rgba(249, 250, 251, 0.5)' }}
                className="hover:bg-gray-50 transition-colors duration-150"
              >
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center">
                    <span className="text-sm font-medium text-gray-900">{code.code}</span>
                    <button
                      onClick={() => copyToClipboard(code.code)}
                      className="ml-2 text-gray-400 hover:text-gray-600"
                    >
                      <ClipboardDocumentIcon className="h-4 w-4" />
                    </button>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      code.status === 'enabled'
                        ? 'bg-green-100 text-green-800'
                        : 'bg-red-100 text-red-800'
                    }`}
                  >
                    {code.status === 'enabled' ? (
                      <CheckCircleIcon className="mr-1 h-4 w-4" />
                    ) : (
                      <XCircleIcon className="mr-1 h-4 w-4" />
                    )}
                    {code.status === 'enabled' ? '启用' : '禁用'}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {new Date(code.createdAt).toLocaleString('zh-CN')}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {new Date(code.expiresAt).toLocaleString('zh-CN')}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  <div className="flex items-center">
                    <div className="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">
                      <div
                        className="h-full bg-indigo-600 rounded-full"
                        style={{
                          width: `${((code.accounts?.length || 0) / code.maxAccounts) * 100}%`,
                        }}
                      />
                    </div>
                    <span className="ml-2">
                      {code.accounts?.length || 0}/{code.maxAccounts}
                    </span>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <button
                    onClick={() => viewDetails(code)}
                    className="text-indigo-600 hover:text-indigo-900 mr-4 inline-flex items-center"
                  >
                    <EyeIcon className="h-4 w-4 mr-1" />
                    查看详情
                  </button>
                  <button
                    onClick={() => toggleStatus(code)}
                    className={`inline-flex items-center ${
                      code.status === 'enabled'
                        ? 'text-red-600 hover:text-red-900'
                        : 'text-green-600 hover:text-green-900'
                    }`}
                  >
                    {code.status === 'enabled' ? (
                      <XCircleIcon className="h-4 w-4 mr-1" />
                    ) : (
                      <CheckCircleIcon className="h-4 w-4 mr-1" />
                    )}
                    {code.status === 'enabled' ? '禁用' : '启用'}
                  </button>
                </td>
              </motion.tr>
            ))}
          </tbody>
        </table>
      </div>
      {renderPagination()}
    </div>
  );

  return (
    <div className="min-h-screen bg-gray-50">
      <Toaster position="top-right" />
      <nav className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <motion.h1 
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-indigo-500 to-purple-500"
              >
                Cursor Reset 管理系统
              </motion.h1>
            </div>
            <div className="flex items-center">
              <button
                onClick={handleLogout}
                className="inline-flex items-center px-4 py-2 text-sm text-gray-700 hover:text-gray-900"
              >
                <ArrowLeftOnRectangleIcon className="w-5 h-5 mr-2" />
                退出登录
              </button>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-white shadow-sm rounded-lg mb-6 overflow-hidden"
        >
          <div className="px-6 py-4 border-b border-gray-200">
            <h2 className="text-lg font-medium text-gray-900">生成激活码</h2>
          </div>
          <form onSubmit={generateActivationCode} className="p-6">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  有效期
                </label>
                <select
                  value={newCode.duration}
                  onChange={(e) => setNewCode({ ...newCode, duration: Number(e.target.value) })}
                  className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 rounded-lg transition-colors duration-200"
                  disabled={isGenerating}
                >
                  <option value={5}>5天</option>
                  <option value={15}>15天</option>
                  <option value={30}>30天</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  可激活账户数
                </label>
                <input
                  type="number"
                  min="1"
                  max="100"
                  value={newCode.maxAccounts}
                  onChange={(e) => setNewCode({ ...newCode, maxAccounts: Number(e.target.value) })}
                  required
                  disabled={isGenerating}
                  className="mt-1 block w-full border-gray-300 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                />
              </div>
              <div className="flex items-end">
                <button
                  type="submit"
                  disabled={isGenerating || isLoading}
                  className="w-full inline-flex justify-center items-center px-4 py-2 border border-transparent text-sm font-medium rounded-lg text-white bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200 transform hover:scale-[1.02]"
                >
                  {isGenerating ? (
                    <>
                      <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      生成中...
                    </>
                  ) : (
                    <>
                      <PlusIcon className="w-5 h-5 mr-2" />
                      生成激活码
                    </>
                  )}
                </button>
              </div>
            </div>
          </form>
        </motion.div>

        {isLoading ? (
          <div className="flex justify-center items-center py-12">
            <svg className="animate-spin h-10 w-10 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span className="ml-3 text-lg text-gray-700">加载中...</span>
          </div>
        ) : (
          renderTable()
        )}

        <AnimatePresence>
          {showModal && selectedCode && (
            <Dialog
              as={motion.div}
              static
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              open={showModal}
              onClose={() => setShowModal(false)}
              className="fixed inset-0 z-10 overflow-y-auto"
            >
              <div className="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
                <Dialog.Overlay
                  as={motion.div}
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  exit={{ opacity: 0 }}
                  className="fixed inset-0 bg-gray-500 bg-opacity-40 transition-opacity"
                />

                <motion.div
                  initial={{ opacity: 0, scale: 0.95 }}
                  animate={{ opacity: 1, scale: 1 }}
                  exit={{ opacity: 0, scale: 0.95 }}
                  className="inline-block align-bottom bg-white rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-2xl sm:w-full sm:p-6"
                >
                  <div className="sm:flex sm:items-start">
                    <div className="mt-3 text-center sm:mt-0 sm:text-left w-full">
                      <Dialog.Title as="h3" className="text-lg leading-6 font-medium text-gray-900">
                        激活码详情
                      </Dialog.Title>
                      <div className="mt-4">
                        <div className="flex items-center justify-between">
                          <p className="text-sm text-gray-500">激活码：{selectedCode.code}</p>
                          <button
                            onClick={() => copyToClipboard(selectedCode.code)}
                            className="text-gray-400 hover:text-gray-600"
                          >
                            <ClipboardDocumentIcon className="h-5 w-5" />
                          </button>
                        </div>
                        {renderAccountsList()}
                      </div>
                    </div>
                  </div>
                  <div className="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
                    <button
                      type="button"
                      className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:w-auto sm:text-sm"
                      onClick={() => setShowModal(false)}
                    >
                      关闭
                    </button>
                  </div>
                </motion.div>
              </div>
            </Dialog>
          )}
        </AnimatePresence>
      </main>
    </div>
  );
} 