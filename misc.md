# Links

- https://peter.bourgon.org/blog/2016/07/11/context.html  How are users of MyMiddleware supposed to know what its dependencies are?
- https://npf.io/2014/04/mocking-functions-in-go/
- https://github.com/Smarp/funcmock
- Dependency injection in go https://habr.com/company/funcorp/blog/372199/
  - https://blog.drewolson.org/dependency-injection-in-go/

```
Господи, какой ужас. Простая, явная и прозрачная инициализация сервиса заменяется на неявную магическую хрень, абсолютно 
не читаемую и подверженную ошибкам при исполнении.

Все таки Java головного мозга очень тяжело поддается лечению.
```

# ApiStack

org.apache.commons.logging.LogFactory by org.apache.commons.logging.impl.LogFactoryImpl
com.triniforce.server.srvapi.IDatabaseInfo by com.triniforce.db.test.DBTestCase$2
com.triniforce.server.srvapi.IIdDef by com.triniforce.server.plugins.kernel.IdDef
--- ApiStack entry:
com.projectkaiser.wfp.api.IWFPApi by com.projectkaiser.wfp.api.WFPApi
com.triniforce.server.srvapi.filesystem.IFileApi$IMeta by com.triniforce.server.plugins.kernel.fileapi.Meta
com.triniforce.server.srvapi.IIdGenerator by com.triniforce.server.plugins.kernel.IdGenerator
com.triniforce.server.srvapi.license.ILicenseInfo by com.triniforce.server.plugins.kernel.ext.api.LicenseInfo
com.triniforce.server.srvapi.services.IServices by com.triniforce.server.plugins.kernel.services.Services
com.triniforce.server.srvapi.i18n.IServerBundleCache by com.triniforce.server.srvapi.i18n.ServerBundleCache
com.triniforce.server.plugins.kernel.ep.external_classes.IExternalClasses by com.triniforce.server.plugins.kernel.ep.external_classes.ExternalClasses
com.triniforce.server.srvapi.IBlobIdGenerator by com.triniforce.server.plugins.kernel.IdGenerator
com.triniforce.server.srvapi.IVersionedBlob by com.triniforce.server.plugins.kernel.filesystem.VersionedBlob
com.triniforce.server.srvapi.filesystem.IFiletypeManager by com.triniforce.db.test.ServerTestCase$TestServer
com.projectkaiser.calendar.ICalendarManager by com.projectkaiser.calendar.CalendarManager
com.triniforce.server.srvapi.filesystem.IFileApi$IGeneral by com.triniforce.server.plugins.kernel.fileapi.General
com.triniforce.server.srvapi.auth.IDbUser by com.triniforce.server.plugins.kernel.files.FTUser$DbUser
com.triniforce.server.srvapi.auth.ICurrentUser by com.triniforce.server.plugins.kernel.CurrentUser
com.triniforce.server.srvapi.ITaskExecutors by com.triniforce.server.plugins.kernel.TaskExecutors
com.triniforce.server.srvapi.INamedDbId by com.triniforce.server.plugins.kernel.tables.TNamedDbId
com.triniforce.server.srvapi.IProjectStructureVersions by com.triniforce.server.plugins.kernel.ext.api.ProjectStructureVersions
org.apache.commons.logging.LogFactory by org.apache.commons.logging.impl.LogFactoryImpl
com.triniforce.server.plugins.kernel.ext.api.IMailer by com.triniforce.server.plugins.kernel.ext.api.Mailer
java.util.Locale by java.util.Locale
com.triniforce.server.srvapi.IMessageBus by com.triniforce.server.plugins.kernel.ext.messagebus.PKMessageBus
com.triniforce.server.srvapi.IServerMode by com.triniforce.server.plugins.kernel.BasicServerCorePlugin$1
com.triniforce.qsync.intf.IQSyncManager by com.triniforce.qsync.impl.QSyncManager
com.triniforce.war.api.IBasicDiag by com.triniforce.war.BasicDiag
com.triniforce.server.srvapi.ITimedLock2 by com.triniforce.server.plugins.kernel.TimedLock2
com.triniforce.server.srvapi.ac.IAccessResolvers by com.triniforce.server.plugins.kernel.ac.AccessResolvers
com.triniforce.server.srvapi.flatviews.IFlatViewFileSystemConnector by com.triniforce.server.plugins.kernel.flatviews.FlatViewFileSystemConnector
com.triniforce.server.srvapi.flatviews.IQueuedCreator by com.triniforce.server.plugins.kernel.flatviews.QueuedCreator
com.triniforce.server.srvapi.filesystem.IFileTriggers by com.triniforce.server.plugins.kernel.filesystem.FileTriggers
com.triniforce.server.srvapi.IServerParameters by com.triniforce.server.srvapi.SrvApiEmu
com.projectkaiser.custom_fields.ICustomFields by com.projectkaiser.custom_fields.APICustomFields
com.triniforce.extensions.IPKExtensionPoint by com.triniforce.db.test.ServerTestCase$TestServer
com.triniforce.server.srvapi.auth.IAuthenticators by com.triniforce.server.plugins.kernel.Authenticators
com.triniforce.server.srvapi.ISODbInfo by com.triniforce.db.test.ServerTestCase$TestServer
java.util.TimeZone by sun.util.calendar.ZoneInfo
com.triniforce.server.srvapi.IMiscIdGenerator by com.triniforce.server.plugins.kernel.MiscIdGenerator
com.triniforce.utils.IProfilerStack by com.triniforce.utils.Profiler$ProfilerStack
com.projectkaiser.wfp.api.IWFPCache by com.projectkaiser.wfp.api.WFPCache
com.triniforce.server.srvapi.auth.IMsgDigests by com.triniforce.server.plugins.kernel.MsgDigests
com.triniforce.server.srvapi.ISrvSmartTranFactory by com.triniforce.server.plugins.kernel.SrvSmartTranFactory
com.triniforce.server.srvapi.ITreeCache by com.triniforce.server.plugins.kernel.filesystem.TreeCache
com.triniforce.server.srvapi.IBlobData by com.triniforce.server.plugins.kernel.filesystem.BlobData
com.triniforce.server.srvapi.IServer by com.triniforce.db.test.ServerTestCase$TestServer
com.triniforce.server.srvapi.IDatabaseInfo by com.triniforce.server.plugins.kernel.BasicServer$1
com.triniforce.server.srvapi.IDbQueueFactory by com.triniforce.server.plugins.kernel.DbQueueFactory
com.triniforce.utils.IProfiler by com.triniforce.utils.Profiler
com.triniforce.server.srvapi.IBasicServer by com.triniforce.db.test.ServerTestCase$TestServer
com.triniforce.server.srvapi.flatviews.IFlatViewsManager by com.triniforce.server.plugins.kernel.flatviews.FlatViewsManager
com.triniforce.server.srvapi.auth.ISessions by com.triniforce.server.plugins.kernel.Sessions
com.triniforce.server.srvapi.filesystem.IFileApi$IDiag by com.triniforce.server.plugins.kernel.fileapi.Diag
com.triniforce.server.srvapi.IOlRawResponseBuilder by com.triniforce.server.plugins.kernel.outline.OlRawResponseBuilder
com.triniforce.server.srvapi.filesystem.IFileApi$IErase by com.triniforce.server.plugins.kernel.fileapi.Erase
com.triniforce.utils.ITime by com.triniforce.db.test.BasicServerApiEmu
com.projectkaiser.auth.ICachedGoogleTokenValidator by com.projectkaiser.auth.CachedGoogleTokenValidator
com.triniforce.server.srvapi.ISrvPrepSqlGetter by com.triniforce.server.plugins.kernel.SrvPrepSqlGetter
com.triniforce.server.srvapi.ISOQuery by com.triniforce.db.test.ServerTestCase$TestServer
com.triniforce.server.srvapi.flatviews.IIndexStorage by com.triniforce.server.plugins.kernel.flatviews.IndexStorage
com.triniforce.server.srvapi.IFlatCache by com.triniforce.server.plugins.kernel.filesystem.FlatCache
com.triniforce.server.srvapi.auth.IPersistentSessions by com.triniforce.server.plugins.kernel.PersistentSessions
com.triniforce.server.srvapi.IUserProperty by com.triniforce.server.plugins.kernel.filesystem.UserProperty
com.triniforce.server.srvapi.IPooledConnection by com.triniforce.db.test.BasicServerTestCase$Pool
com.triniforce.server.srvapi.IIdDef by com.triniforce.server.plugins.kernel.IdDef
com.triniforce.server.plugins.kernel.ext.api.PTRecurringTasks by com.triniforce.server.plugins.kernel.ext.api.PTRecurringTasks
com.triniforce.server.srvapi.IThrdWatcherRegistrator by com.triniforce.server.plugins.kernel.ThrdWatcherRegistrator
com.triniforce.server.srvapi.filesystem.IFileApi$IAccessControl by com.triniforce.server.plugins.kernel.fileapi.AccessControl
com.triniforce.server.srvapi.i18n.IServerBundle by com.triniforce.server.srvapi.i18n.ServerBundle
--- ApiStack entry:
java.sql.Connection by org.apache.commons.dbcp.PoolingDataSource$PoolGuardConnectionWrapper
com.triniforce.server.srvapi.ISrvSmartTran by com.triniforce.server.plugins.kernel.SrvSmartTran
--- ApiStack entry:
com.triniforce.server.srvapi.ISrvSmartTranExtenders$IRefCountHashMap by com.triniforce.server.plugins.kernel.BasicServerCorePlugin$RefCountMapTrnExtender$RefCountMap
--- ApiStack entry:
com.triniforce.server.srvapi.ISrvSmartTranExtenders$IFiniter by com.triniforce.server.plugins.kernel.BasicServerCorePlugin$FiniterExtender$Finiter
--- ApiStack entry:
com.triniforce.server.srvapi.ISrvSmartTranExtenders$ILocker by com.triniforce.server.plugins.kernel.BasicServerCorePlugin$LockerExtender$Locker
--- ApiStack entry:
com.triniforce.server.srvapi.ITransactionWriteLock2 by com.triniforce.server.plugins.kernel.BasicServerCorePlugin$TransactionWriteLock
--- ApiStack entry:
com.triniforce.server.srvapi.filesystem.IFileDataFactory by com.triniforce.server.plugins.kernel.filesystem.FileDataFactory
com.triniforce.server.srvapi.filesystem.IFileModificationFactory by com.triniforce.server.plugins.kernel.filesystem.FileModificationFactory
--- ApiStack eof
