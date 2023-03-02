# coding=utf-8
# *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import copy
import warnings
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
from .. import _utilities

import types

__config__ = pulumi.Config('twingate')


class _ExportableConfig(types.ModuleType):
    @property
    def api_token(self) -> Optional[str]:
        """
        The access key for API operations. You can retrieve this from the Twingate Admin Console
        ([documentation](https://docs.twingate.com/docs/api-overview)). Alternatively, this can be specified using the
        TWINGATE_API_TOKEN environment variable.
        """
        return __config__.get('apiToken')

    @property
    def http_max_retry(self) -> int:
        """
        Specifies a retry limit for the http requests made. The default value is 10. Alternatively, this can be specified using
        the TWINGATE_HTTP_MAX_RETRY environment variable
        """
        return __config__.get_int('httpMaxRetry') or 5

    @property
    def http_timeout(self) -> int:
        """
        Specifies a time limit in seconds for the http requests made. The default value is 10 seconds. Alternatively, this can
        be specified using the TWINGATE_HTTP_TIMEOUT environment variable
        """
        return __config__.get_int('httpTimeout') or 10

    @property
    def network(self) -> Optional[str]:
        """
        Your Twingate network ID for API operations. You can find it in the Admin Console URL, for example:
        `autoco.twingate.com`, where `autoco` is your network ID Alternatively, this can be specified using the TWINGATE_NETWORK
        environment variable.
        """
        return __config__.get('network')

    @property
    def url(self) -> Optional[str]:
        """
        The default is 'twingate.com' This is optional and shouldn't be changed under normal circumstances.
        """
        return __config__.get('url')
